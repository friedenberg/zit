package user_ops

import (
	"io"
	"os"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/charlie/open_file_guard"
	"github.com/friedenberg/zit/hotel/organize_text"
)

type ReadOrganizeFile struct {
	io.Reader
}

func (c ReadOrganizeFile) RunWithFile(p string) (ot organize_text.Text, err error) {
	var f *os.File

	if f, err = open_file_guard.Open(p); err != nil {
		err = errors.Error(err)
		return
	}

	defer open_file_guard.Close(f)

	c.Reader = f
	ot, err = c.Run()
	c.Reader = nil

	return
}

func (c ReadOrganizeFile) Run() (ot organize_text.Text, err error) {
	ot = organize_text.NewEmpty()

	_, err = ot.ReadFrom(c.Reader)

	return
}
