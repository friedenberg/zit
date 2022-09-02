package zettel_named

type NamedFilter interface {
	IncludeNamedZettel(Zettel) bool
}

type FilterAll struct{}

func (f FilterAll) IncludeNamedZettel(_ Zettel) bool {
	return true
}
