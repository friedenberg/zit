package collections_coding

import "github.com/friedenberg/zit/src/bravo/collections"

type EncoderLike[T any] interface {
	Encode(*T) (int64, error)
}

func EncoderToWriter[T any](el EncoderLike[T]) collections.WriterFunc[*T] {
	return func(e *T) (err error) {
		_, err = el.Encode(e)
		return
	}
}
