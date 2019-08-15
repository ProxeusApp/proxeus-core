package validate

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/c2h5oh/datasize"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
)

var ftRepository FTypes

func init() {
	//add pdf to the document validation
	pdfType := filetype.GetType("pdf")
	matchers.Document[pdfType] = func(buf []byte) bool {
		return filetype.IsType(buf, pdfType)
	}
	//provide a static repository to make visible what types are supported
	vagueImg := make([]types.Type, len(matchers.Image))
	collectTypes(vagueImg, matchers.Image)
	vagueVid := make([]types.Type, len(matchers.Video))
	collectTypes(vagueVid, matchers.Video)
	vagueAudio := make([]types.Type, len(matchers.Audio))
	collectTypes(vagueAudio, matchers.Audio)
	vagueDoc := make([]types.Type, len(matchers.Document))
	collectTypes(vagueDoc, matchers.Document)
	vagueArchive := make([]types.Type, len(matchers.Archive))
	collectTypes(vagueArchive, matchers.Archive)
	vagueFont := make([]types.Type, len(matchers.Font))
	collectTypes(vagueFont, matchers.Font)
	ftRepository = FTypes{
		Vague: map[string][]types.Type{
			vagueTypeImage:    vagueImg,
			vagueTypeVideo:    vagueVid,
			vagueTypeAudio:    vagueAudio,
			vagueTypeDocument: vagueDoc,
			vagueTypeArchive:  vagueArchive,
			vagueTypeFont:     vagueFont,
		},
		Exact: filetype.Types,
	}
}

func collectTypes(col []types.Type, from matchers.Map) {
	i := 0
	for k := range from {
		col[i] = k
		i++
	}
}

const vagueTypeImage = "image"
const vagueTypeVideo = "video"
const vagueTypeAudio = "audio"
const vagueTypeDocument = "document"
const vagueTypeArchive = "archive"
const vagueTypeFont = "font"

type FTypes struct {
	Vague map[string][]types.Type
	Exact map[string]types.Type
}

type FileType struct {
	Exact bool   `json:"exact"` //vague                  or exact
	Kind  string `json:"kind"`  //(image,video,audio...) or (jpg,png,mp4,mp3...)
}

type TmpFile struct {
	*os.File
}

func (me *TmpFile) Close() error {
	if me.File != nil {
		err := me.File.Close()
		_ = os.Remove(me.Name())
		return err
	}
	return nil
}

func (me *TmpFile) CloseWithoutRemove() error {
	return me.File.Close()
}

func (r Rules) FileType() *FileType {
	if ftif, ok := r["file"]; ok && ftif != nil {
		bts, err := json.Marshal(ftif)
		if err != nil {
			return nil
		}
		var ft *FileType
		_ = json.Unmarshal(bts, &ft)
		return ft
	}
	return nil
}

func FileTypes() FTypes {
	return ftRepository
}

func FileWithExt(src io.Reader, fileExtension, min, max string) (*TmpFile, error) {
	return FileWithOptions(src, FileType{Exact: true, Kind: fileExtension}, min, max)
}

func FileWithType(src io.Reader, ft types.Type, min, max string) (*TmpFile, error) {
	return FileWithOptions(src, FileType{Exact: true, Kind: ft.Extension}, min, max)
}

func FileWithOptions(src io.Reader, ft FileType, min, max string) (*TmpFile, error) {
	rules := Rules{}
	if ft.Kind != "" {
		rules["file"] = Rules{
			"exact": ft.Exact,
			"kind":  ft.Kind,
		}
	}
	if min != "" {
		rules["min"] = min
	}
	if max != "" {
		rules["max"] = max
	}
	return File(src, rules)
}

func File(src io.Reader, rules Rules) (*TmpFile, error) {
	validr := newValidator(nil, rules)
	if validr == nil {
		return nil, os.ErrInvalid
	}
	var err error
	var tmpFile *TmpFile
	defer func() {
		if err != nil && tmpFile != nil {
			//ensure we cleanup if any error occurs by closing and removing the tmp file
			_ = tmpFile.Close()
		}
	}()
	//read FileType, minSize and maxSize from rules
	var tmpf *os.File
	tmpf, err = ioutil.TempFile("", "validate")
	if err != nil {
		return nil, err
	}
	tmpFile = &TmpFile{tmpf}

	var minBytes int64 = -1
	var maxBytes int64 = -1
	ft := rules.FileType()
	if v, ok := validr.hasStrValueFor("min"); ok {
		min, err := hummanReadableStringToUInt64Bytes(v)
		if err != nil {
			validr.add(&Error{Msg: msgBadDefinitionOfMin})
			err = *validr.errs
			return nil, *validr.errs
		}
		minBytes = int64(min)
	}
	if v, ok := validr.hasStrValueFor("max"); ok {
		max, err := hummanReadableStringToUInt64Bytes(v)
		if err != nil {
			validr.add(&Error{Msg: msgBadDefinitionOfMax})
			err = *validr.errs
			return nil, err
		}
		maxBytes = int64(max)
	}
	_, err = validateStream(tmpFile, src, ft, minBytes, maxBytes)
	if err != nil {
		validr.add(&Error{Msg: err.Error()})
		err = *validr.errs
		return nil, *validr.errs
	}
	err = tmpFile.CloseWithoutRemove()
	if err != nil {
		return nil, err
	}
	tmpf, err = os.Open(tmpFile.Name())
	if err != nil {
		return nil, err
	}
	tmpFile = &TmpFile{tmpf}
	return tmpFile, nil
}

func (me *validator) file() {

}

func (me *validator) fileStream(dst io.Writer, src io.Reader) {
}

//validateStream breaks if there is an error in the type, min or max validation
//if any of the options are undefined, the validation will pass
//no error is returned for undefined!
//
//ft nil for undefined
//minBytes -1 for undefined
//maxBytes -1 for undefined
func validateStream(dst io.Writer, src io.Reader, ft *FileType, minBytes, maxBytes int64) (written int64, err error) {
	buf := make([]byte, 32*1024)
	needToValidate := true
	if ft == nil {
		needToValidate = false
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			if needToValidate {
				//validate the type by reading the beginning of the content
				if nr >= 262 {
					err = validateType(buf, ft)
					if err != nil {
						return
					}
					//validation succeeded, set flag to false as we don't need to check the rest of the content
					needToValidate = false
				} else {
					err = ErrFileSizeToLowForValidation
					return
				}
			}

			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
				if maxBytes != -1 && written > maxBytes {
					err = ErrFileSizeToHigh
					return
				}
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	if minBytes != -1 && written < minBytes {
		err = ErrFileSizeToLow
		return
	}
	if maxBytes != -1 && written > maxBytes {
		err = ErrFileSizeToHigh
		return
	}
	return written, err
}

var ErrFileSizeToLowForValidation = fmt.Errorf("file size to low for validation")
var ErrFileSizeToLow = fmt.Errorf("file size to low")
var ErrFileSizeToHigh = fmt.Errorf("file size to high")
var ErrWrongFileType = fmt.Errorf("wrong file type")
var ErrUnknownFileType = fmt.Errorf("define a supported type")

func validateType(buf []byte, ft *FileType) error {
	if ft == nil || ft.Kind == "" {
		//no type definition no validation -> no error
		return nil
	}
	t := strings.ToLower(ft.Kind)
	if ft.Exact {
		tt := filetype.GetType(t)
		if tt == filetype.Unknown {
			return ErrUnknownFileType
		}
		if filetype.IsType(buf, tt) {
			return nil
		}
	} else {
		switch t {
		case vagueTypeImage:
			if filetype.IsImage(buf) {
				return nil
			}
		case vagueTypeVideo:
			if filetype.IsVideo(buf) {
				return nil
			}
		case vagueTypeAudio:
			if filetype.IsAudio(buf) {
				return nil
			}
		case vagueTypeDocument:
			if filetype.IsDocument(buf) {
				return nil
			}
		case vagueTypeArchive:
			if filetype.IsArchive(buf) {
				return nil
			}
		case vagueTypeFont:
			if filetype.IsFont(buf) {
				return nil
			}
		default:
			return ErrUnknownFileType
		}
	}
	return ErrWrongFileType
}

func hummanReadableStringToUInt64Bytes(hummanReadable string) (uint64, error) {
	var v datasize.ByteSize
	err := v.UnmarshalText([]byte(hummanReadable))
	if err != nil {
		return 0, err
	}
	return v.Bytes(), nil
}
