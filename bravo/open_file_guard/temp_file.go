package open_file_guard

import (
	"io/ioutil"
	"os"
)

func TempFile() (f *os.File, err error) {
	openFilesGuardInstance.Lock()

	if f, err = ioutil.TempFile("", ""); err != nil {
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
