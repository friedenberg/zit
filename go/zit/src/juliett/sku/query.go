package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type (
	Queryable interface {
		// AppendMatchToQueryPath(*Transacted, *QueryPath) error
		ContainsSku(TransactedGetter) bool
	}

	Query interface {
		Queryable
		interfaces.Stringer
		// Each(schnittstellen.FuncIter[Query]) error
	}

	SigilGetter interface {
		GetSigil() ids.Sigil
	}

	QueryWithSigilAndObjectId interface {
		Query
		SigilGetter
		ContainsObjectId(*ids.ObjectId) bool
	}

	// Used by store_verzeichnisse.binary*
	PrimitiveQueryGroup interface {
		Get(genres.Genre) (QueryWithSigilAndObjectId, bool)
		SigilGetter
		HasHidden() bool
	}

	FuncQuery = func(
		QueryGroup,
		interfaces.FuncIter[*Transacted],
	) (err error)

	FuncPrimitiveQuery = func(
		PrimitiveQueryGroup,
		interfaces.FuncIter[*Transacted],
	) (err error)

	QueryGroup interface {
		PrimitiveQueryGroup
		Query
		SetIncludeHistory()
		GetTags() ids.TagSet
		GetTypes() ids.TypeSet
	}

	ExternalQueryOptions struct {
		ids.RepoId
		ExcludeUntracked  bool
		ExcludeRecognized bool
	}

	ExternalQuery struct {
		ids.RepoId
		QueryGroup
		ExternalQueryOptions
	}
)

func MakePrimitiveQueryGroup() PrimitiveQueryGroup {
	return MakePrimitiveQueryGroupWithSigils(ids.SigilHistory, ids.SigilHidden)
}

func MakePrimitiveQueryGroupWithSigils(ss ...ids.Sigil) PrimitiveQueryGroup {
	return &primitiveQueryGroup{Sigil: ids.MakeSigil(ss...)}
}

type primitiveQueryGroup struct {
	ids.Sigil
}

func (qg *primitiveQueryGroup) SetIncludeHistory() {
	qg.Add(ids.SigilHistory)
}

func (qg *primitiveQueryGroup) HasHidden() bool {
	return false
}

func (qg *primitiveQueryGroup) Get(_ genres.Genre) (QueryWithSigilAndObjectId, bool) {
	return qg, true
}

func (s *primitiveQueryGroup) ContainsSku(_ TransactedGetter) bool {
	return true
}

func (s *primitiveQueryGroup) ContainsObjectId(_ *ids.ObjectId) bool {
	return false
}

func (s *primitiveQueryGroup) GetSigil() ids.Sigil {
	return s.Sigil
}

func (s *primitiveQueryGroup) Each(_ interfaces.FuncIter[Query]) error {
	return nil
}
