package query

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type (
	Matcher        = sku.QueryBase
	MatcherSigil   = sku.Query
	MatchableAdder = sku.MatchableAdder
)

type Cwd interface {
	Matcher
	GetCwdFDs() fd.Set
	GetKennungForFD(*fd.FD) (*kennung.Kennung2, error)
}

type ImplicitEtikettenGetter interface {
	GetImplicitEtiketten(*kennung.Etikett) kennung.EtikettSet
}

func VisitAllMatchers(
	f schnittstellen.FuncIter[Matcher],
	matchers ...Matcher,
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
			func(m Matcher) (err error) {
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

func MakeCheckedOutQueryFunc(
	m sku.QueryBase,
) schnittstellen.FuncIter[*sku.CheckedOut] {
	if m == nil {
		return collections.MakeWriterNoop[*sku.CheckedOut]()
	}

	return func(col *sku.CheckedOut) (err error) {
		if !m.ContainsMatchable(&col.External.Transacted) {
			err = iter.MakeErrStopIteration()
			return
		}

		return
	}
}

type (
	FuncReaderTransactedLikePtr func(schnittstellen.FuncIter[*sku.Transacted]) error
	FuncQueryTransactedLikePtr  func(*Group, schnittstellen.FuncIter[*sku.Transacted]) error
)
