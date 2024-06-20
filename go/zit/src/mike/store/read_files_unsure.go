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
	existing sku.CheckedOutMutableSet,
) error

func (s *Store) QueryUnsure(
	qg *query.Group,
	o UnsureMatchOptions,
	f IterMatching,
) (err error) {
	selbstMetadateiSansTaiToZettels := make(
		map[sha.Bytes]sku.CheckedOutMutableSet,
		s.GetCwdFiles().Len(),
	)

	bezToZettels := make(
		map[string]sku.CheckedOutMutableSet,
		s.GetCwdFiles().Len(),
	)

	var l sync.Mutex

	if err = s.ReadFilesUnsure(
		qg,
		func(co *sku.CheckedOut) (err error) {
			sh := &co.External.Metadatei.Shas.SelbstMetadateiSansTai

			if sh.IsNull() {
				return
			}

			l.Lock()
			defer l.Unlock()

			clone := sku.GetCheckedOutPool().Get()
			sku.CheckedOutResetter.ResetWith(clone, co)

			if o.Contains(UnsureMatchTypeMetadateiSansTaiHistory) {
				k := sh.GetBytes()
				existing, ok := selbstMetadateiSansTaiToZettels[k]

				if !ok {
					existing = sku.MakeCheckedOutMutableSet()
				}

				if err = existing.Add(clone); err != nil {
					err = errors.Wrap(err)
					return
				}

				selbstMetadateiSansTaiToZettels[k] = existing
			}

			if o.Contains(UnsureMatchTypeBezeichnung) {
				k := co.External.Metadatei.Bezeichnung.String()
				existing, ok := bezToZettels[k]

				if !ok {
					existing = sku.MakeCheckedOutMutableSet()
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
