package zettel_named

import "io"

type NamedFilter interface {
	IncludeNamedZettel(*Zettel) bool
}

type FilterAll struct{}

func (f FilterAll) IncludeNamedZettel(_ *Zettel) bool {
	return true
}

type WriterFilter struct {
	NamedFilter
}

func (w WriterFilter) WriteZettelNamed(z *Zettel) (err error) {
	if !w.IncludeNamedZettel(z) {
		err = io.EOF
	}

	return
}
