package commands

import (
	"flag"
	"os"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/charlie/files"
	"code.linenisgreat.com/zit/src/kilo/organize_text"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
	"code.linenisgreat.com/zit/src/papa/user_ops"
)

type FormatOrganize struct {
	Flags organize_text.Flags
}

func init() {
	registerCommand(
		"format-organize",
		func(f *flag.FlagSet) Command {
			c := &FormatOrganize{
				Flags: organize_text.MakeFlags(),
			}

			c.Flags.AddToFlagSet(f)

			return c
		},
	)
}

func (c *FormatOrganize) Run(u *umwelt.Umwelt, args ...string) (err error) {
	c.Flags.Konfig = u.Konfig()

	if len(args) != 1 {
		err = errors.Errorf("expected exactly one input argument")
		return
	}

	var f *os.File

	if f, err = files.Open(args[0]); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	var ot *organize_text.Text

	readOrganizeTextOp := user_ops.ReadOrganizeFile{
		Umwelt: u,
		Reader: f,
	}

	if ot, err = readOrganizeTextOp.Run(nil); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO move Abbr as required arg
	ot.Options = c.Flags.GetOptions(
		u.Konfig().PrintOptions,
		nil,
		u.SkuFormatOldOrganize(),
		u.SkuFmtNewOrganize(),
		u.MakeKennungExpanders(),
	)

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
