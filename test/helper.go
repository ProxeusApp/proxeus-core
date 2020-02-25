package test

import (
	"os"
	"path"
	"runtime"
)

// This is used to always run tests from the project's root folder instead of from the package.
// Just import the (unused) package _ "github.com/ProxeusApp/test" and call this function in tests where a disk read is required
func init() {
	_, filename, _, _ := runtime.Caller(0)
	rootFolder := path.Join(path.Dir(filename), "..")

	err := os.Chdir(rootFolder)
	if err != nil {
		panic(err)
	}
}
