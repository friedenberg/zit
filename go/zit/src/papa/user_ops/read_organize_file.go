package user_ops

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/kilo/organize_text"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type ReadOrganizeFile struct {
	*umwelt.Umwelt
	io.Reader
}

func (c ReadOrganizeFile) RunWithFile(
	p string,
	q matcher.Query,
) (ot *organize_text.Text, err error) {
	var f *os.File

	if f, err = files.Open(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer files.Close(f)

	c.Reader = f
	ot, err = c.Run(q)
	c.Reader = nil

	return
}

func (c ReadOrganizeFile) Run(q matcher.Query) (ot *organize_text.Text, err error) {
	otFlags := organize_text.MakeFlags()
	c.Umwelt.ApplyToOrganizeOptions(&otFlags.Options)

	if ot, err = organize_text.New(
		otFlags.GetOptions(
			c.Umwelt.Konfig().PrintOptions,
			q,
			c.SkuFormatOldOrganize(),
			c.SkuFmtNewOrganize(),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = ot.ReadFrom(c.Reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = ot.Refine(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}