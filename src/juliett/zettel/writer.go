package zettel

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func MakeWriterZettel(
	wf schnittstellen.FuncIter[*Objekte],
) schnittstellen.FuncIter[*Transacted] {
	return func(z *Transacted) (err error) {
		if err = wf(&z.Objekte); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}
