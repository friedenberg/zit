package organize_text

import (
	"github.com/friedenberg/zit/src/alfa/bezeichnung"
	"github.com/friedenberg/zit/src/alfa/errors"
)

type newZettel struct {
	bezeichnung.Bezeichnung
}

func (z *newZettel) Set(v string) (err error) {
	remaining := v

	if remaining[:2] != "- " {
		err = errors.Errorf("expected '- ', but got '%s'", remaining[:2])
		return
	}

	remaining = remaining[2:]

	if err = z.Bezeichnung.Set(remaining); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
