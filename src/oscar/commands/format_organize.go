package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/india/organize_text"
	"github.com/friedenberg/zit/src/mike/umwelt"
	"github.com/friedenberg/zit/src/november/user_ops"
)

type FormatOrganize struct {
	organize_text.Options
}

func init() {
	registerCommand(
		"format-organize",
		func(f *flag.FlagSet) Command {
			c := &FormatOrganize{
				Options: organize_text.MakeOptions(),
			}

			c.Options.AddToFlagSet(f)

			return c
		},
	)
}

func (c *FormatOrganize) Run(u *umwelt.Umwelt, args ...string) (err error) {
	c.Options.Konfig = u.Konfig()

	if len(args) != 1 {
		err = errors.Errorf("expected exactly one input argument")
		return
	}

	var f *os.File

	if f, err = files.Open(args[0]); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer files.Close(f)

	var ot *organize_text.Text

	readOrganizeTextOp := user_ops.ReadOrganizeFile{
		Umwelt:  u,
		Reader:  f,
	}

	if ot, err = readOrganizeTextOp.Run(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ot.Options = c.Options

	if err = ot.Refine(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = ot.WriteTo(os.Stdout); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
