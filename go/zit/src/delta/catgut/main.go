package catgut

import "code.linenisgreat.com/zit/src/alfa/schnittstellen"

type (
	StringFormatReader[T any] interface {
		ReadStringFormat(*RingBuffer, T) (int64, error)
	}

	StringFormatReadWriter[T any] interface {
		StringFormatReader[T]
		schnittstellen.StringFormatWriter[T]
	}
)
