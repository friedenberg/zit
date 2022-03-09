package user_ops

import (
	"fmt"
	"os"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/golf/organize_text"
)

type ReadOrganizeFile struct {
}

func (c ReadOrganizeFile) Run(p string) (ot organize_text.Text, err error) {
	for {
		var f *os.File

		if f, err = open_file_guard.Open(p); err != nil {
			err = _Error(err)
			return
		}

		defer open_file_guard.Close(f)

		ot = organize_text.NewEmpty()

		if _, err = ot.ReadFrom(f); err != nil {
			//TODO move to parent
			if c.handleReadChangesError(err) {
				continue
			} else {
				stdprinter.Errf("aborting organize\n")
				return
			}
		}

		break
	}

	return
}

func (co ReadOrganizeFile) handleReadChangesError(err error) (tryAgain bool) {
	var errorRead organize_text.ErrorRead

	if err != nil && !errors.As(err, &errorRead) {
		stdprinter.Errf("unrecoverable organize read failure: %s", err)
		tryAgain = false
		return
	}

	stdprinter.Errf("reading changes failed: %q\n", err)
	stdprinter.Errf("would you like to edit and try again? (y/*)\n")

	var answer rune
	var n int

	if n, err = fmt.Scanf("%c", &answer); err != nil {
		tryAgain = false
		stdprinter.Errf("failed to read answer: %s", err)
		return
	}

	if n != 1 {
		tryAgain = false
		stdprinter.Errf("failed to read at exactly 1 answer: %s", err)
		return
	}

	if answer == 'y' || answer == 'Y' {
		tryAgain = true
	}

	return
}
