package organize_text

import (
	"sort"

	"github.com/friedenberg/zit/src/hotel/zettel_named"
)

type zettelSet map[zettel]bool

func makeZettelZetFromSetNamed(set zettel_named.SetNamed) (zs zettelSet) {
	zs = makeZettelSet()

	for _, z := range set {
		zs.Add(zettel{Hinweis: z.Hinweis.String(), Bezeichnung: z.Stored.Zettel.Bezeichnung.String()})
	}

	return
}

func makeZettelSet() zettelSet {
	return make(map[zettel]bool)
}

func (zs *zettelSet) Add(z zettel) {
	(*zs)[z] = true
}

func (zs *zettelSet) Del(z zettel) {
	delete(*zs, z)
}

func (zs zettelSet) sorted() (sorted []zettel) {
	sorted = make([]zettel, len(zs))
	i := 0

	for z, _ := range zs {
		sorted[i] = z
		i++
	}

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Bezeichnung == sorted[j].Bezeichnung {
			return sorted[i].Hinweis < sorted[j].Hinweis
		} else {
			return sorted[i].Bezeichnung < sorted[j].Bezeichnung
		}
	})

	return
}

func (zs zettelSet) Contains(z zettel) bool {
	_, ok := zs[z]
	return ok
}

func (a zettelSet) Equals(b zettelSet) bool {
	if len(a) != len(b) {
		return false
	}

	for z, _ := range a {
		if !b.Contains(z) {
			return false
		}
	}

	return true
}

func (a zettelSet) Copy() (b zettelSet) {
	b = makeZettelSet()

	for z, _ := range a {
		b[z] = true
	}

	return
}
