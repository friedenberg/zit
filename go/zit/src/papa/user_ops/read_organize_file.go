package user_ops

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type ReadOrganizeFile struct{}

func (c ReadOrganizeFile) RunWithPath(
	u *local_working_copy.Repo,
	p string,
	repoId ids.RepoId,
) (ot *organize_text.Text, err error) {
	var f *os.File

	if f, err = files.Open(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	if ot, err = c.Run(
		u,
		f,
		organize_text.NewMetadata(repoId),
	); err != nil {
		err = errors.Wrapf(err, "Path: %q", p)
		return
	}

	return
}

func (c ReadOrganizeFile) Run(
	u *local_working_copy.Repo,
	r io.Reader,
	om organize_text.Metadata,
) (ot *organize_text.Text, err error) {
	otFlags := organize_text.MakeFlags()
	u.ApplyToOrganizeOptions(&otFlags.Options)

	o := otFlags.GetOptionsWithMetadata(
		u.GetConfig().GetCLIConfig().PrintOptions,
		u.SkuFormatBoxCheckedOutNoColor(),
		u.GetStore().GetAbbrStore().GetAbbr(),
		sku.ObjectFactory{},
		om,
	)

	if ot, err = organize_text.New(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = ot.ReadFrom(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
