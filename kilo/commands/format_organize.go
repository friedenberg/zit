package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/golf/organize_text"
	"github.com/friedenberg/zit/juliett/user_ops"
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

func (c *FormatOrganize) Run(u _Umwelt, args ...string) (err error) {
	if len(args) != 1 {
		err = errors.Errorf("expected exactly one input argument")
		return
	}

	// stdoutIsTty := open_file_guard.IsTty(os.Stdout)
	// stdinIsTty := open_file_guard.IsTty(os.Stdin)

	// if !stdinIsTty && !stdoutIsTty {
	// 	logz.Print("neither stdin or stdout is a tty")
	// 	logz.Print("generate organize, read from stdin, commit")

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

	if _, err = ot.WriteTo(os.Stdout); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
