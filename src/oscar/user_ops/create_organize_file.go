package user_ops

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/organize_text"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type CreateOrganizeFile struct {
	*umwelt.Umwelt
	organize_text.Options
}

func (c CreateOrganizeFile) RunAndWrite(
	zettels zettel_transacted.MutableSet,
	w io.WriteCloser,
) (results *organize_text.Text, err error) {
	defer errors.Deferred(&err, w.Close)

	if results, err = c.Run(zettels); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = results.WriteTo(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c CreateOrganizeFile) Run(zettels zettel_transacted.MutableSet) (results *organize_text.Text, err error) {
	if results, err = organize_text.New(c.Options); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
