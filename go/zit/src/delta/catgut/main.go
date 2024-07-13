package catgut

import "code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"

type (
	StringFormatReader[T any] interface {
		ReadStringFormat(*RingBuffer, T) (int64, error)
	}

	StringFormatReadWriter[T any] interface {
		StringFormatReader[T]
		interfaces.StringFormatWriter[T]
	}
)
