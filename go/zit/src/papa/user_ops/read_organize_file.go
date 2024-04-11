package user_ops

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/charlie/files"
	"code.linenisgreat.com/zit/src/india/query"
	"code.linenisgreat.com/zit/src/kilo/organize_text"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type ReadOrganizeFile struct {
	*umwelt.Umwelt
	io.Reader
}

func (c ReadOrganizeFile) RunWithFile(
	p string,
	q *query.Group,
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

func (c ReadOrganizeFile) Run(q *query.Group) (ot *organize_text.Text, err error) {
	otFlags := organize_text.MakeFlags()
	c.ApplyToOrganizeOptions(&otFlags.Options)

	if ot, err = organize_text.New(
		otFlags.GetOptions(
			c.Umwelt.Konfig().PrintOptions,
			q,
			c.SkuFormatOldOrganize(),
			c.SkuFmtNewOrganize(),
			c.MakeKennungExpanders(),
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
