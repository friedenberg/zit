package iter

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func AddString[T schnittstellen.Value[T]](
	sb *strings.Builder,
) schnittstellen.FuncIter[T] {
	return func(e T) (err error) {
		sb.WriteString(e.String())
		sb.WriteString(" ")

		return
	}
}

func MakeFuncIterNoOp[T any]() schnittstellen.FuncIter[T] {
	return func(e T) (err error) {
		return
	}
}
