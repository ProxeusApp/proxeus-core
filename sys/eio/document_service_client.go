package eio

import (
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
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

//Compile packs the provided files into a ZIP and sends it to the document-service to be compiled as the format you have provided.
//if format is empty, it will take the default one which is PDF.
func (ds *DocumentServiceClient) Compile(template Template) (resp *http.Response, err error) {
	var odtTmplFile *os.File
	odtTmplFile, err = os.Open(template.TemplatePath)
	if err != nil {
		return
	}
	fstat, er := odtTmplFile.Stat()
	if er != nil {
		err = er
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
			if odtTmplFile != nil {
				odtTmplFile.Close()
			}
			if requestWriter != nil {
				requestWriter.Close()
			}
		}()
		header, err := zip.FileInfoHeader(fstat)
		if err != nil {
			return
		}
		header.Name = "tmpl.odt"

		// Change to deflate to gain better compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate
		odtZipWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			return
		}
		_, err = io.CopyN(odtZipWriter, odtTmplFile, fstat.Size())
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
				assetFile, err := os.Open(assetPath)
				if err != nil {
					return
				}
				defer assetFile.Close()
				assetStat, err := assetFile.Stat()
				if err != nil {
					return
				}

				assetHeader, err := zip.FileInfoHeader(assetStat)
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
				_, err = io.CopyN(assetWriter, assetFile, assetStat.Size())
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
	client.Timeout = time.Duration(time.Second * 20)
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

//Vars returns a list of the vars contained in the provided template.
//You can filter them with the prefix if needed.
func (ds *DocumentServiceClient) Vars(tmplPath, prefix string) ([]string, error) {
	odtTmplFile, err := os.Open(tmplPath)
	if err != nil {
		return nil, err
	}
	fstat, err := odtTmplFile.Stat()
	if err != nil {
		return nil, err
	}
	requestReader, requestWriter := io.Pipe()
	go func() {
		defer func() {
			if odtTmplFile != nil {
				_ = odtTmplFile.Close()
			}
			if requestWriter != nil {
				_ = requestWriter.Close()
			}
		}()
		_, err = io.CopyN(requestWriter, odtTmplFile, fstat.Size())
		if err != nil {
			return
		}
	}()
	req, er := http.NewRequest("POST", ds.makeUrl("vars"), requestReader)
	if er != nil {
		return nil, er
	}
	if prefix != "" {
		q := req.URL.Query()
		q.Add("prefix", prefix)
		req.URL.RawQuery = q.Encode()
	}
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/zip")

	client := &http.Client{}
	client.Timeout = time.Duration(time.Second * 20)
	resp, er := client.Do(req)
	if er != nil {
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

//Download the template assistance extension for your writer.
func (ds *DocumentServiceClient) DownloadExtension(arch string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", ds.makeUrl("extension"), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("os", OSArchType(arch).String())
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	client.Timeout = time.Duration(time.Second * 10)
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

//Download the how it works document.
func (ds *DocumentServiceClient) DownloadHowItWorksPDF(writer io.Writer) error {
	return ds.downloadPDF("how-it-works", false, writer)
}

//Download the example document.
func (ds *DocumentServiceClient) DownloadExamplePDF(raw bool, writer io.Writer) error {
	return ds.downloadPDF("example", raw, writer)
}

func (ds *DocumentServiceClient) downloadPDF(method string, raw bool, writer io.Writer) error {
	rp := ""
	if raw {
		rp = "?raw"
	}
	req, err := http.NewRequest("GET", ds.makeUrl(method)+rp, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	client.Timeout = time.Duration(time.Second * 10)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return os.ErrInvalid
	}
	err = checkGzip(resp)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, resp.Body)
	resp.Body.Close()
	return err
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
