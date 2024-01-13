package iter

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func AddStringToBuilder[T schnittstellen.Value[T]](
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

// func MakeFuncIter[IN any, OUT any](m Matcher) schnittstellen.FuncIter[T] {
// 	return func(e T) (err error) {
// 		if !m.ContainsMatchable(e) {
// 			err = iter.MakeErrStopIteration()
// 			return
// 		}

// 		return
// 	}
// }
