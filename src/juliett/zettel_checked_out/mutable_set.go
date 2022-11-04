package zettel_checked_out

import (
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/golf/zettel_external"
)

type MutableSet struct {
	collections.MutableSet[*Zettel]
}

func MakeMutableSetUnique(c int) MutableSet {
	return MutableSet{
		MutableSet: collections.MakeMutableSet(
			func(sz *Zettel) string {
				if sz == nil {
					return ""
				}

				return collections.MakeKey(
					sz.Internal.Kopf,
					sz.Internal.Mutter,
					sz.Internal.Schwanz,
					sz.Internal.Named.Hinweis,
					sz.Internal.Named.Stored.Sha,
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

				return collections.MakeKey(sz.Internal.Named.Hinweis)
			},
		),
	}
}

func (s MutableSet) ToSliceZettelsExternal() (out []zettel_external.Zettel) {
	out = make([]zettel_external.Zettel, 0, s.Len())

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
			p := z.External.ZettelFD.Path

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
