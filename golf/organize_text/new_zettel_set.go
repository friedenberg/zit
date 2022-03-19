package organize_text

type newZettelSet map[newZettel]bool

func makeNewZettelSet() newZettelSet {
	return make(map[newZettel]bool)
}

func (zs newZettelSet) Add(z newZettel) {
	zs[z] = true
}

func (zs newZettelSet) Del(z newZettel) {
	delete(zs, z)
}

func (zs newZettelSet) Contains(z newZettel) bool {
	_, ok := zs[z]
	return ok
}
