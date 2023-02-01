package files

import (
	"io/ioutil"
	"os"
)

func TempDir() (d string, err error) {
	openFilesGuardInstance.Lock()

	if d, err = ioutil.TempDir("", ""); err != nil {
		openFilesGuardInstance.Unlock()
	}

	return
}

func TempFile(d string) (f *os.File, err error) {
	openFilesGuardInstance.Lock()

	if f, err = ioutil.TempFile(d, ""); err != nil {
		openFilesGuardInstance.Unlock()
	}

	return
}

func TempFileInDir(dir string) (f *os.File, err error) {
	openFilesGuardInstance.Lock()

	if f, err = ioutil.TempFile(dir, ""); err != nil {
		openFilesGuardInstance.Unlock()
	}

	return
}

func TempFileWithPattern(pattern string) (f *os.File, err error) {
	openFilesGuardInstance.Lock()

	if f, err = ioutil.TempFile("", pattern); err != nil {
		openFilesGuardInstance.Unlock()
	}

	return
}
