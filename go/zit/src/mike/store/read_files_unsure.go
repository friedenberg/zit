package store

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

type UnsureMatchType byte

const (
	UnsureMatchTypeNone = UnsureMatchType(iota << 1)
	UnsureMatchTypeMetadateiSansTaiHistory
	UnsureMatchTypeBezeichnung

	UnsureMatchTypeAll = UnsureMatchType(^byte(0))
)

func (a UnsureMatchType) Contains(b UnsureMatchType) bool {
	return a&b == b
}

type UnsureMatchOptions struct {
	UnsureMatchType
	Filter schnittstellen.FuncIter[*sku.Transacted]
}

func UnsureMatchOptionsDefault() UnsureMatchOptions {
	return UnsureMatchOptions{
		UnsureMatchType: UnsureMatchTypeMetadateiSansTaiHistory | UnsureMatchTypeBezeichnung,
	}
}

type IterMatching func(
	mt UnsureMatchType,
	sk *sku.Transacted,
	existing sku.CheckedOutLikeMutableSet,
) error

func (s *Store) QueryUnsure(
	qg *query.Group,
	o UnsureMatchOptions,
	f IterMatching,
) (err error) {
	selbstMetadateiSansTaiToZettels := make(
		map[sha.Bytes]sku.CheckedOutLikeMutableSet,
		s.GetCwdFiles().Len(),
	)

	bezToZettels := make(
		map[string]sku.CheckedOutLikeMutableSet,
		s.GetCwdFiles().Len(),
	)

	var l sync.Mutex

	if err = s.cwdFiles.QueryUnsure(
		qg,
		func(col sku.CheckedOutLike) (err error) {
			e := col.GetSkuExternalLike().GetSku()
			sh := &e.Metadatei.Shas.SelbstMetadateiSansTai

			if sh.IsNull() {
				return
			}

			l.Lock()
			defer l.Unlock()

			clone := col.Clone()

			if o.Contains(UnsureMatchTypeMetadateiSansTaiHistory) {
				k := sh.GetBytes()
				existing, ok := selbstMetadateiSansTaiToZettels[k]

				if !ok {
					existing = sku.MakeCheckedOutLikeMutableSet()
				}

				if err = existing.Add(clone); err != nil {
					err = errors.Wrap(err)
					return
				}

				selbstMetadateiSansTaiToZettels[k] = existing
			}

			if o.Contains(UnsureMatchTypeBezeichnung) {
				k := e.Metadatei.Bezeichnung.String()
				existing, ok := bezToZettels[k]

				if !ok {
					existing = sku.MakeCheckedOutLikeMutableSet()
				}

				if err = existing.Add(clone); err != nil {
					err = errors.Wrap(err)
					return
				}

				bezToZettels[k] = existing
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO create a new query group for all of history
	qg.SetIncludeHistory()

	if len(selbstMetadateiSansTaiToZettels) > 0 || len(bezToZettels) > 0 {
		if err = s.QueryWithoutCwd(
			qg,
			func(sk *sku.Transacted) (err error) {
				if o.Filter != nil {
					if err = o.Filter(sk); err != nil {
						err = errors.Wrap(err)
						return
					}
				}

				if err = sk.CalculateObjekteShas(); err != nil {
					err = errors.Wrap(err)
					return
				}

				sh := &sk.Metadatei.Shas.SelbstMetadateiSansTai

				if sh.IsNull() {
					return
				}

				{
					k := sh.GetBytes()
					existing, ok := selbstMetadateiSansTaiToZettels[k]

					if !ok {
						return
					}

					if err = f(
						UnsureMatchTypeMetadateiSansTaiHistory,
						sk,
						existing,
					); err != nil {
						err = errors.Wrap(err)
						return
					}
				}

				{
					k := sk.Metadatei.Bezeichnung.String()
					existing, ok := bezToZettels[k]

					if !ok {
						return
					}

					if err = f(
						UnsureMatchTypeBezeichnung,
						sk,
						existing,
					); err != nil {
						err = errors.Wrap(err)
						return
					}
				}

				return
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
