package user_ops

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/read_write_repo_local"
)

// TODO support using query results for organize population
type CreateOrganizeFile struct {
	*read_write_repo_local.Repo
	organize_text.Options
}

func (c CreateOrganizeFile) RunAndWrite(
	w io.Writer,
) (results *organize_text.Text, err error) {
	if results, err = c.Run(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = results.WriteTo(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c CreateOrganizeFile) Run() (results *organize_text.Text, err error) {
	if results, err = organize_text.New(c.Options); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
