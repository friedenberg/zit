package matcher_proto

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/delta/gattungen"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type QueryGroup interface {
	Get(g gattung.Gattung) (s MatcherSigil, ok bool)
	GetCwdFDs() fd.Set
	GetExplicitCwdFDs() fd.Set
	GetEtiketten() kennung.EtikettSet
	GetTypen() kennung.TypSet
	GetGattungen() gattungen.Set
	Matcher
}

type QueryGroupBuilder interface {
	BuildQueryGroup(...string) (QueryGroup, error)
}

type Matcher interface {
	ContainsMatchable(*sku.Transacted) bool
	schnittstellen.Stringer
	Each(schnittstellen.FuncIter[Matcher]) error
}

type MatcherSigil interface {
	Matcher
	GetSigil() kennung.Sigil
}

type MatcherSigilPtr interface {
	MatcherSigil
	AddSigil(kennung.Sigil)
}

type MatcherParentPtr interface {
	Matcher
	Add(Matcher) error
}

type MatchableAdder interface {
	AddMatchable(*sku.Transacted) error
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

func SplitGattungenByHistory(qg QueryGroup) (schwanz, all kennung.Gattung) {
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

func MakeFilterFromQuery(
	ms QueryGroup,
) schnittstellen.FuncIter[*sku.CheckedOut] {
	if ms == nil {
		return collections.MakeWriterNoop[*sku.CheckedOut]()
	}

	return func(col *sku.CheckedOut) (err error) {
		if !ms.ContainsMatchable(&col.External.Transacted) {
			err = iter.MakeErrStopIteration()
			return
		}

		return
	}
}

type (
	FuncReaderTransactedLikePtr func(schnittstellen.FuncIter[*sku.Transacted]) error
	// FuncSigilTransactedLikePtr  func(MatcherSigil, schnittstellen.FuncIter[*sku.Transacted]) error
	FuncQueryTransactedLikePtr func(QueryGroup, schnittstellen.FuncIter[*sku.Transacted]) error
)

func MakeFuncReaderTransactedLikePtr(
	ms QueryGroup,
	fq FuncQueryTransactedLikePtr,
) FuncReaderTransactedLikePtr {
	return func(f schnittstellen.FuncIter[*sku.Transacted]) (err error) {
		return fq(ms, f)
	}
}
