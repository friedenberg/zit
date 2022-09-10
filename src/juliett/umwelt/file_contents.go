package umwelt

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/friedenberg/zit/src/bravo/files"
)

func (u Umwelt) ContentsOfZitFile(p ...string) (contents string, err error) {
	basePath := u.DirZit()

	var f *os.File

	p = append([]string{basePath}, p...)

	if f, err = files.Open(path.Join(p...)); err != nil {
		return
	}

	defer files.Close(f)

	var b []byte

	if b, err = ioutil.ReadAll(f); err != nil {
		return
	}

	contents = string(b)

	return
}
