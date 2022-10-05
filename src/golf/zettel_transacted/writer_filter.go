package zettel_transacted

type WriterFilter struct {
	predicate func(Zettel) bool
	inner     Writer
}

func MakeWriterFilter(w Writer, f func(Zettel) bool) *WriterFilter {
	return &WriterFilter{
		predicate: f,
		inner:     w,
	}
}

func (w WriterFilter) WriteZettelTransacted(z Zettel) (err error) {
	if !w.predicate(z) {
		return
	}

	return w.inner.WriteZettelTransacted(z)
}
