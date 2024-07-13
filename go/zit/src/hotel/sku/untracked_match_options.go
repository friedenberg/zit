package sku

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
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

func (a UnsureMatchType) MakeMatchMap() UnsureMatchMaps {
	maps := UnsureMatchMaps{
		Lookup: make(map[UnsureMatchType]UnsureMatchMap),
	}

	if a.Contains(UnsureMatchTypeMetadateiSansTaiHistory) {
		maps.Lookup[UnsureMatchTypeMetadateiSansTaiHistory] = UnsureMatchMap{
			UnsureMatchType: UnsureMatchTypeMetadateiSansTaiHistory,
			Lookup:          make(map[sha.Bytes]CheckedOutLikeMutableSet),
		}
	}

	if a.Contains(UnsureMatchTypeBezeichnung) {
		maps.Lookup[UnsureMatchTypeBezeichnung] = UnsureMatchMap{
			UnsureMatchType: UnsureMatchTypeBezeichnung,
			Lookup:          make(map[sha.Bytes]CheckedOutLikeMutableSet),
		}
	}

	return maps
}

type UnsureMatchOptions struct {
	UnsureMatchType
}

func UnsureMatchOptionsDefault() UnsureMatchOptions {
	return UnsureMatchOptions{
		UnsureMatchType: UnsureMatchTypeMetadateiSansTaiHistory | UnsureMatchTypeBezeichnung,
	}
}

type IterMatching func(
	mt UnsureMatchType,
	sk *Transacted,
	existing CheckedOutLikeMutableSet,
) error

type UnsureMatchMap struct {
	UnsureMatchType
	Lookup map[sha.Bytes]CheckedOutLikeMutableSet
}

type UnsureMatchMaps struct {
	Lookup map[UnsureMatchType]UnsureMatchMap
}

func (umm UnsureMatchMaps) Len() int {
	l := 0

	for _, v := range umm.Lookup {
		l += len(v.Lookup)
	}

	return l
}

func MakeUnsureMatchMapsCollector(
	umm UnsureMatchMaps,
) interfaces.FuncIter[CheckedOutLike] {
	var l sync.Mutex

	return func(col CheckedOutLike) (err error) {
		e := col.GetSkuExternalLike().GetSku()

		l.Lock()
		defer l.Unlock()

		clone := col.Clone()

		for t, v := range umm.Lookup {
			var k sha.Bytes

			switch t {
			case UnsureMatchTypeMetadateiSansTaiHistory:
				k = e.Metadatei.Shas.SelbstMetadateiSansTai.GetBytes()

			case UnsureMatchTypeBezeichnung:
				k = sha.FromString(e.Metadatei.Bezeichnung.String()).GetBytes()

			default:
				continue
			}

			existing, ok := v.Lookup[k]

			if !ok {
				existing = MakeCheckedOutLikeMutableSet()
			}

			if err = existing.Add(clone); err != nil {
				err = errors.Wrap(err)
				return
			}

			v.Lookup[k] = existing
		}

		return
	}
}

func MakeUnsureMatchMapsMatcher(
	umm UnsureMatchMaps,
	f IterMatching,
) interfaces.FuncIter[*Transacted] {
	return func(sk *Transacted) (err error) {
		for t, v := range umm.Lookup {
			var k sha.Bytes

			switch t {
			case UnsureMatchTypeMetadateiSansTaiHistory:
				k = sk.Metadatei.Shas.SelbstMetadateiSansTai.GetBytes()

			case UnsureMatchTypeBezeichnung:
				k = sha.FromString(sk.Metadatei.Bezeichnung.String()).GetBytes()

			default:
				continue
			}

			existing, ok := v.Lookup[k]

			if !ok {
				continue
			}

			if err = f(
				t,
				sk,
				existing,
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		return
	}
}
