package eio

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys/file"
)

type (
	Format                string
	OSArchType            string
	DocumentServiceClient struct{ Url string }
)

const (
	PDF  Format = "pdf"
	ODT  Format = "odt"
	DOCX Format = "docx"
	DOC  Format = "doc"

	Linux  OSArchType = "linux_x86_64"
	Win64  OSArchType = "win_x86_64"
	Win32  OSArchType = "win_x86"
	Darwin OSArchType = "mac_x86_64"
)

type Template struct {
	Format       Format      //the result format, default PDF
	Data         interface{} //the data the template is going to resolved with
	TemplatePath string      //odt template
	Assets       []string    //assets that are referenced in Data and you want to include into the ZIP
	EmbedError   bool        //print compilation error in the document
}

// Compile packs the provided files into a ZIP and sends it to the document-service to be compiled as the format you have provided.
// if format is empty, it will take the default one which is PDF.
func (ds *DocumentServiceClient) Compile(db storage.FilesIF, template Template) (resp *http.Response, err error) {
	var templateBuf bytes.Buffer
	err = db.Read(template.TemplatePath, &templateBuf)
	if err != nil {
		return
	}
	requestReader, requestWriter := io.Pipe()
	go func() {
		zipWriter := zip.NewWriter(requestWriter)
		defer func() {
			if zipWriter != nil && zipWriter.Close() != nil {
				requestWriter.Close()
				return
			}
			if requestWriter != nil {
				requestWriter.Close()
			}
		}()
		header, err := zip.FileInfoHeader(
			file.InMemoryFileInfo{Path: "tmpl.odt", Len: templateBuf.Len()})
		if err != nil {
			return
		}

		// Change to deflate to gain better compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate
		odtZipWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			return
		}
		_, err = io.Copy(odtZipWriter, &templateBuf)
		if err != nil {
			return
		}

		//insert data as a json file in the ZIP
		jsonBytes, err := json.Marshal(template.Data)
		if err != nil {
			jsonBytes = []byte("{}")
		}

		djh := &zip.FileHeader{
			Name:   "data.json",
			Method: zip.Deflate,
		}

		dataJsonWriter, err := zipWriter.CreateHeader(djh)
		if err != nil {
			return
		}
		_, err = dataJsonWriter.Write(jsonBytes)
		if err != nil {
			return
		}
		//insert assets
		if len(template.Assets) > 0 {
			assetToZIPWriter := func(assetPath string, zipWriter *zip.Writer) {
				var asssetBuf bytes.Buffer
				err := db.Read(assetPath, &asssetBuf)
				if err != nil {
					return
				}
				assetHeader, err := zip.FileInfoHeader(file.InMemoryFileInfo{
					Path: filepath.Base(assetPath),
					Len:  asssetBuf.Len(),
				})
				if err != nil {
					return
				}

				// Change to deflate to gain better compression
				// see http://golang.org/pkg/archive/zip/#pkg-constants
				assetHeader.Method = zip.Deflate
				assetWriter, err := zipWriter.CreateHeader(assetHeader)
				if err != nil {
					return
				}
				_, err = io.Copy(assetWriter, &asssetBuf)
				if err != nil {
					return
				}
			}
			for _, assetPath := range template.Assets {
				assetToZIPWriter(assetPath, zipWriter)
			}
		}
	}()
	req, er := http.NewRequest("POST", ds.makeUrl("compile"), requestReader)
	if er != nil {
		err = er
		return
	}
	if template.Format != "" {
		q := req.URL.Query()
		q.Add("format", template.Format.String())
		req.URL.RawQuery = q.Encode()
	}
	if template.EmbedError {
		q := req.URL.Query()
		//just need to provide the parameter
		//no value needed
		q.Add("error", "1")
		req.URL.RawQuery = q.Encode()
	}
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/zip")

	client := &http.Client{}
	client.Timeout = time.Second * 20
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	err = checkGzip(resp)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		var b []byte
		defer resp.Body.Close()
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		err = errors.New(string(b))
	}
	return
}

// Vars returns a list of the vars contained in the provided template.
// You can filter them with the prefix if needed.
func (ds *DocumentServiceClient) Vars(templateBuf *bytes.Buffer) ([]string, error) {
	requestReader, requestWriter := io.Pipe()
	go func() {
		io.Copy(requestWriter, templateBuf)
		requestWriter.Close()
		requestReader.Close()
	}()
	req, er := http.NewRequest("POST", ds.makeUrl("vars"), requestReader)
	if er != nil {
		return nil, er
	}
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/zip")

	client := &http.Client{}
	client.Timeout = time.Second * 20
	resp, err := client.Do(req)
	if err != nil {
		return nil, er
	}
	err = checkGzip(resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		bts, er := ioutil.ReadAll(resp.Body)
		if er != nil {
			return nil, er
		}
		var vars []string
		err = json.Unmarshal(bts, &vars)
		if err != nil {
			return nil, err
		}
		return vars, nil
	} else {
		var b []byte
		defer resp.Body.Close()
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = errors.New(string(b))
		return nil, err
	}
}

// Download the template assistance extension for your writer.
func (ds *DocumentServiceClient) DownloadExtension(arch string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", ds.makeUrl("extension"), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("os", OSArchType(arch).String())
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	client.Timeout = time.Second * 10
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = os.ErrInvalid
		return
	}
	err = checkGzip(resp)
	if err != nil {
		return
	}
	return
}

func checkGzip(resp *http.Response) (err error) {
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		resp.Body, err = gzip.NewReader(resp.Body)
	}
	return
}

func (ds *DocumentServiceClient) makeUrl(method string) string {
	u, err := url.Parse(ds.Url)
	if err != nil {
		return ""
	}
	u.Path = path.Join(u.Path, method)
	return u.String()
}

func (f Format) String() string {
	switch f {
	case ODT:
		return string(ODT)
	case DOCX:
		return string(DOCX)
	case DOC:
		return string(DOC)
	default:
		return string(PDF)
	}
}

func (f OSArchType) String() string {
	switch f {
	case Win64:
		return string(Win64)
	case Win32:
		return string(Win32)
	case Darwin:
		return string(Darwin)
	default:
		return string(Linux)
	}
}
