package zettel

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/india/konfig"
)

func MakeWriterKonfig(
	k konfig.Compiled,
) schnittstellen.FuncIter[*Transacted] {
	errors.TodoP1("switch to sigils")
	errors.TodoP3("add efficient parsing of hiding tags")

	if k.IncludeHidden {
		return collections.MakeWriterNoop[*Transacted]()
	}

	return func(z *Transacted) (err error) {
		if err = z.GetMetadatei().Etiketten.Each(
			func(e kennung.Etikett) (err error) {
				p := e.String()

				for _, t := range k.EtikettenHiddenStringsSlice {
					if strings.HasPrefix(p, t) {
						err = collections.MakeErrStopIteration()
						return
					}
				}

				return
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		t := k.GetApproximatedTyp(z.GetTyp()).ApproximatedOrActual()

		if t != nil && t.Akte.Archived {
			err = collections.MakeErrStopIteration()
			return
		}

		return
	}
}
