package zettel

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/typ_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/india/transacted"
	"github.com/friedenberg/zit/src/kilo/konfig"
)

func MakeWriterKonfig(
	k konfig.Compiled,
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) schnittstellen.FuncIter[*transacted.Zettel] {
	errors.TodoP1("switch to sigils")
	errors.TodoP3("add efficient parsing of hiding tags")

	if k.IncludeHidden {
		return collections.MakeWriterNoop[*transacted.Zettel]()
	}

	return func(z *transacted.Zettel) (err error) {
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

		var ta *typ_akte.V0

		if ta, err = tagp.GetAkte(t.GetAkteSha()); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer tagp.PutAkte(ta)

		if t != nil && ta.Archived {
			err = collections.MakeErrStopIteration()
			return
		}

		return
	}
}