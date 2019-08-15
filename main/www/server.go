package www

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"path/filepath"
)

type MyServer struct {
	quit chan os.Signal
}

func (ms *MyServer) Close() {
	if ms.quit != nil {
		ms.quit <- os.Interrupt
	}
}

type MyHTMLTemplateLoader struct {
	BaseDir  string
	MoreDirs *[]string
}

// Abs calculates the path to a given template. Whenever a path must be resolved
// due to an import from another template, the base equals the parent template's path.
func (htl *MyHTMLTemplateLoader) Abs(base, name string) (absPath string) {
	if filepath.IsAbs(name) {
		return name
	}

	// Our own base dir has always priority; if there's none
	// we use the path provided in base.
	var err error
	if htl.BaseDir == "" {
		if base == "" {
			base, err = os.Getwd()
			if err != nil {
				panic(err)
			}
			absPath = filepath.Join(base, name)
			htl.checkPath(&name, &absPath)
			return
		}
		absPath = filepath.Join(filepath.Dir(base), name)
		htl.checkPath(&name, &absPath)
		return
	}
	absPath = filepath.Join(htl.BaseDir, name)
	htl.checkPath(&name, &absPath)
	return
}

func (htl *MyHTMLTemplateLoader) checkPath(relPath, absPath *string) {
	if htl.MoreDirs != nil && len(*htl.MoreDirs) > 0 {
		var err error
		if _, err = os.Stat(*absPath); err == nil {
			return
		}
		newPath := ""
		for _, dirPath := range *htl.MoreDirs {
			newPath = filepath.Join(dirPath, *relPath)
			if _, err = os.Stat(newPath); err == nil {
				*absPath = newPath
				break
			}
		}
	}
}

// Get returns an io.Reader where the template's content can be read from.
func (htl *MyHTMLTemplateLoader) Get(path string) (io.Reader, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buf), nil
}
