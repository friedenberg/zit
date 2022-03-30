package user_ops

import (
	"os"

	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/golf/organize_text"
)

type ReadOrganizeFile struct {
}

func (c ReadOrganizeFile) Run(p string) (ot organize_text.Text, err error) {
	var f *os.File

	if f, err = open_file_guard.Open(p); err != nil {
		err = _Error(err)
		return
	}

	defer open_file_guard.Close(f)

	ot = organize_text.NewEmpty()

	_, err = ot.ReadFrom(f)

	return
}
