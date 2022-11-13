package zettel_named

type FilterAnd []NamedFilter

func (f FilterAnd) IncludeNamedZettel(z *Zettel) bool {
	for _, f1 := range []NamedFilter(f) {
		if !f1.IncludeNamedZettel(z) {
			return false
		}
	}

	return true
}
