package tar

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Given a folder (src), tar-gz' the content and writes it to writer
// Example:
// 	f, _ := os.Create(tmpFile)
//
//	err = Tar(dirToZip, f)
func Tar(src string, writer io.Writer) error {
	gzw := gzip.NewWriter(writer)
	tw := tar.NewWriter(gzw)
	defer func() {
		tw.Close()
		gzw.Close()
	}()

	// walk path
	err := filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}
		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// return on non-regular files (thanks to [kumo](https://medium.com/@komuw/just-like-you-did-fbdd7df829d3) for this suggested update)
		if !fi.Mode().IsRegular() {
			return nil
		}
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}
		f.Close()
		return nil
	})
	return err
}

// Given a .tar.gz file as Reader, decompress the files contained on it to specified destination folder
func Untar(destination string, r io.Reader) (err error) {
	err = os.MkdirAll(destination, 0750)
	if err != nil {
		return
	}
	var gzr *gzip.Reader
	gzr, err = gzip.NewReader(r)
	if err != nil {
		return
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, er := tr.Next()

		switch {
		// if no more files are found return
		case er == io.EOF:
			return
			// return any other error
		case er != nil:
			return
			// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		target := filepath.Join(destination, header.Name)
		// the target location where the dir/file should be created

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if err = os.MkdirAll(target, 0755); err != nil {
				return
			}
			// if it's a file create it
		case tar.TypeReg:
			var f *os.File
			f, err = os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return
			}
			// copy over contents
			if _, err = io.Copy(f, tr); err != nil {
				f.Close()
				return
			}
			f.Close()
		}
	}
}
