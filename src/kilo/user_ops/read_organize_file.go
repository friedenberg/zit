package user_ops

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/open_file_guard"
	"github.com/friedenberg/zit/src/hotel/organize_text"
)

type ReadOrganizeFile struct {
	io.Reader
}

func (c ReadOrganizeFile) RunWithFile(p string) (ot organize_text.Text, err error) {
	var f *os.File

	if f, err = open_file_guard.Open(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer open_file_guard.Close(f)

	c.Reader = f
	ot, err = c.Run()
	c.Reader = nil

	return
}

func (c ReadOrganizeFile) Run() (ot organize_text.Text, err error) {
	if ot, err = organize_text.New(organize_text.Options{}); err != nil {
		err = errors.Wrap(err)
		return
	}

	_, err = ot.ReadFrom(c.Reader)

	return
}
