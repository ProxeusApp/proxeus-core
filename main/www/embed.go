package www

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/labstack/echo"
)

type Embedded struct {
	Asset func(name string) ([]byte, error)
}

func (emb *Embedded) Asset2(name string) ([]byte, error) {
	if name != "" {
		if strings.HasPrefix(name, "/") {
			return emb.Asset(name[1:])
		} else {
			return emb.Asset(name)
		}
	}
	return nil, os.ErrNotExist
}

func (emb *Embedded) FindAssetWithContentType(name string, ct *string) ([]byte, error) {
	if name != "" {
		jj, err := emb.Asset2(name)
		if err != nil {
			return nil, err
		}

		contentTypeMutex.Lock()

		if contentTypeCache == nil {
			contentTypeCache = make(map[string]string)
		}

		contentType := contentTypeCache[name]

		contentTypeMutex.Unlock()

		if contentType == "" {
			li := strings.LastIndex(name, ".")
			if li > 0 {
				ext := name[li:]
				contentType = ctMap[ext]
			}
			if contentType == "" {
				contentType = http.DetectContentType(jj)
			}
			contentTypeMutex.Lock()
			ct2 := contentTypeCache[name]
			if ct2 == "" {
				contentTypeCache[name] = contentType
			}
			contentTypeMutex.Unlock()
		}
		*ct = contentType
		return jj, nil
	}
	return nil, echo.ErrNotFound
}

type EmbeddedTemplateLoader struct {
	Embedded *Embedded
}

// Abs calculates the path to a given template. Whenever a path must be resolved
// due to an import from another template, the base equals the parent template's path.
func (htl *EmbeddedTemplateLoader) Abs(base, name string) (absPath string) {
	if base != "" {
		if htl.exists(name) {
			return name
		}
		name = filepath.Join(filepath.Dir(base), name)
	}
	return name
}

func (htl *EmbeddedTemplateLoader) exists(p string) bool {
	buf, err := htl.Embedded.Asset(p)
	if err != nil {
		return false
	}
	return len(buf) > 0
}

// Get returns an io.Reader where the template's content can be read from.
func (htl *EmbeddedTemplateLoader) Get(path string) (io.Reader, error) {
	if path != "" {
		if htl.Embedded != nil && htl.Embedded.Asset != nil {
			buf, err := htl.Embedded.Asset(path)
			if err != nil {
				return nil, err
			}
			return bytes.NewReader(buf), nil
		}
	}
	return nil, echo.ErrNotFound
}

var (
	contentTypeCache map[string]string
	contentTypeMutex sync.RWMutex

	ctMap = map[string]string{
		".js":    "application/javascript",
		".css":   "text/css",
		".gif":   "image/gif",
		".png":   "image/png",
		".tff":   "application/font-sfnt",
		".otf":   "application/font-sfnt",
		".woff":  "application/font-woff",
		".woff2": "application/font-woff",
		".svg":   "image/svg+xml",
		".eot":   "image/png",
	}
)
