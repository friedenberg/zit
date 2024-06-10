package user_ops

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/charlie/files"
	"code.linenisgreat.com/zit/src/kilo/organize_text"
	"code.linenisgreat.com/zit/src/november/umwelt"
)

type ReadOrganizeFile struct{}

func (c ReadOrganizeFile) RunWithPath(
	u *umwelt.Umwelt,
	p string,
) (ot *organize_text.Text, err error) {
	var f *os.File

	if f, err = files.Open(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	if ot, err = c.Run(u, f); err != nil {
		err = errors.Wrapf(err, "Path: %q", p)
		return
	}

	return
}

func (c ReadOrganizeFile) Run(
	u *umwelt.Umwelt,
	r io.Reader,
) (ot *organize_text.Text, err error) {
	otFlags := organize_text.MakeFlags()
	u.ApplyToOrganizeOptions(&otFlags.Options)

	if ot, err = organize_text.New(
		otFlags.GetOptions(
			u.GetKonfig().PrintOptions,
			nil,
			u.SkuFmtOrganize(),
			u.GetStore().GetAbbrStore().GetAbbr(),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = ot.ReadFrom(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = ot.Refine(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
