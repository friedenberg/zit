package quiter

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func AddStringToBuilder[T interfaces.Value[T]](
	sb *strings.Builder,
) interfaces.FuncIter[T] {
	return func(e T) (err error) {
		sb.WriteString(e.String())
		sb.WriteString(" ")

		return
	}
}

func MakeFuncIterNoOp[T any]() interfaces.FuncIter[T] {
	return func(e T) (err error) {
		return
	}
}
