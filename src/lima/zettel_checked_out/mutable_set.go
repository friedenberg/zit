package zettel_checked_out

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type MutableSet struct {
	schnittstellen.MutableSet[*Zettel]
}

func MakeMutableSetUnique(c int) MutableSet {
	return MutableSet{
		MutableSet: collections.MakeMutableSet(
			func(sz *Zettel) string {
				if sz == nil {
					return ""
				}

				return collections.MakeKey(
					sz.Internal.Sku.Kopf,
					sz.Internal.Sku.Mutter[0],
					sz.Internal.Sku.Mutter[1],
					sz.Internal.Sku.Schwanz,
					sz.Internal.Sku.Kennung,
					sz.Internal.Sku.ObjekteSha,
				)
			},
		),
	}
}

func MakeMutableSetHinweisZettel(c int) MutableSet {
	return MutableSet{
		MutableSet: collections.MakeMutableSet(
			func(sz *Zettel) string {
				if sz == nil {
					return ""
				}

				return collections.MakeKey(sz.Internal.Sku.Kennung)
			},
		),
	}
}

func (s MutableSet) ToSliceZettelsExternal() (out []zettel.External) {
	out = make([]zettel.External, 0, s.Len())

	s.Each(
		func(z *Zettel) (err error) {
			out = append(out, z.External)
			return
		},
	)

	return
}

func (s MutableSet) ToSliceFilesZettelen() (out []string) {
	out = make([]string, 0, s.Len())

	s.Each(
		func(z *Zettel) (err error) {
			p := z.External.FD.Path

			if p != "" {
				out = append(out, p)
			}

			return
		},
	)

	return
}

func (s MutableSet) ToSliceFilesAkten() (out []string) {
	out = make([]string, 0, s.Len())

	s.Each(
		func(z *Zettel) (err error) {
			p := z.External.AkteFD.Path

			if p != "" {
				out = append(out, p)
			}

			return
		},
	)

	return
}
