package query

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type Cwd interface {
	sku.Query
	GetCwdFDs() fd.Set
	GetKennungForFD(*fd.FD) (*kennung.Kennung2, error)
}

type ImplicitEtikettenGetter interface {
	GetImplicitEtiketten(*kennung.Etikett) kennung.EtikettSet
}

func VisitAllMatchers(
	f schnittstellen.FuncIter[sku.Query],
	matchers ...sku.Query,
) (err error) {
	for _, m := range matchers {
		if err = f(m); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		if err = m.Each(
			func(m sku.Query) (err error) {
				return VisitAllMatchers(f, m)
			},
		); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func SplitGattungenByHistory(qg *Group) (schwanz, all kennung.Gattung) {
	err := qg.GetGattungen().Each(
		func(g gattung.Gattung) (err error) {
			m, ok := qg.Get(g)

			if !ok {
				return
			}

			if m.GetSigil().IncludesHistory() {
				all.Add(g)
			} else {
				schwanz.Add(g)
			}

			return
		},
	)

	errors.PanicIfError(err)

	return
}
