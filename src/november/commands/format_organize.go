package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/juliett/organize_text"
	"github.com/friedenberg/zit/src/mike/user_ops"
)

type FormatOrganize struct {
}

func init() {
	registerCommand(
		"format-organize",
		func(f *flag.FlagSet) Command {
			c := &FormatOrganize{}

			return c
		},
	)
}

func (c *FormatOrganize) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) != 1 {
		err = errors.Errorf("expected exactly one input argument")
		return
	}

	var f *os.File

	if f, err = open_file_guard.Open(args[0]); err != nil {
		err = errors.Error(err)
		return
	}

	defer open_file_guard.Close(f)

	var ot organize_text.Text

	readOrganizeTextOp := user_ops.ReadOrganizeFile{
		Reader: f,
	}

	if ot, err = readOrganizeTextOp.Run(); err != nil {
		err = errors.Error(err)
		return
	}

	refiner := organize_text.AssignmentTreeRefiner{
		Enabled:         true,
		UsePrefixJoints: true,
	}

	if err = ot.Refine(refiner); err != nil {
		err = errors.Error(err)
		return
	}

	if _, err = ot.WriteTo(os.Stdout); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
