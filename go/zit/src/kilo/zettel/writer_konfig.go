package zettel

import (
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/delta/typ_akte"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/juliett/konfig"
)

func MakeWriterKonfig(
	k konfig.Compiled,
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) schnittstellen.FuncIter[*sku.Transacted] {
	errors.TodoP1("switch to sigils")
	errors.TodoP3("add efficient parsing of hiding tags")

	if k.IncludeHidden {
		return collections.MakeWriterNoop[*sku.Transacted]()
	}

	return func(z *sku.Transacted) (err error) {
		if err = z.Metadatei.GetEtiketten().Each(
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

		if t == nil {
			log.Err().Printf("empty typ: %s", z.Kennung)
			return
		}

		if tagp == nil {
			panic("typ akte getter putter was nil")
		}

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
