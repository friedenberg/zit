package catgut

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

type (
	StringFormatReader[T any] interface {
		ReadStringFormat(*RingBuffer, T) (int64, error)
	}

	StringFormatReadWriter[T any] interface {
		StringFormatReader[T]
		schnittstellen.StringFormatWriter[T]
	}
)
