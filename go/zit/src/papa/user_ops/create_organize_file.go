package user_ops

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/kilo/organize_text"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

// TODO support using query results for organize population
type CreateOrganizeFile struct {
	*umwelt.Umwelt
	organize_text.Options
}

func (c CreateOrganizeFile) RunAndWrite(
	w io.WriteCloser,
) (results *organize_text.Text, err error) {
	defer errors.DeferredCloser(&err, w)

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
