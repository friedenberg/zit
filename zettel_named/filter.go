package zettel_named

type NamedFilter interface {
	IncludeNamedZettel(Named) bool
}

type FilterAll struct{}

func (f FilterAll) IncludeNamedZettel(_ Named) bool {
	return true
}
