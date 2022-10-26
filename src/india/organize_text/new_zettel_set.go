package organize_text

import "sort"

type newZettelSet map[newZettel]bool

func makeNewZettelSet() newZettelSet {
	return make(map[newZettel]bool)
}

func (zs *newZettelSet) Add(z newZettel) {
	(*zs)[z] = true
}

func (zs *newZettelSet) Del(z newZettel) {
	delete(*zs, z)
}

func (zs newZettelSet) Contains(z newZettel) bool {
	_, ok := zs[z]
	return ok
}

func (zs newZettelSet) sorted() (sorted []newZettel) {
	sorted = make([]newZettel, len(zs))
	i := 0

	for z, _ := range zs {
		sorted[i] = z
		i++
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Bezeichnung < sorted[j].Bezeichnung
	})

	return
}
