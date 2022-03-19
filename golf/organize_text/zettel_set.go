package organize_text

import "sort"

type zettelSet map[zettel]bool

func makeZettelSet() zettelSet {
	return make(map[zettel]bool)
}

func (zs zettelSet) Add(z zettel) {
	zs[z] = true
}

func (zs zettelSet) Del(z zettel) {
	delete(zs, z)
}

func (zs zettelSet) sorted() (sorted []zettel) {
	sorted = make([]zettel, len(zs))
	i := 0

	for z, _ := range zs {
		sorted[i] = z
		i++
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].hinweis < sorted[j].hinweis
	})

	return
}

func (zs zettelSet) Contains(z zettel) bool {
	_, ok := zs[z]
	return ok
}

func (a zettelSet) Copy() (b zettelSet) {
	b = makeZettelSet()

	for z, _ := range a {
		b[z] = true
	}

	return
}
