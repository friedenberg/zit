package zettel_transacted

import (
	"fmt"
	"strings"
)

type SetKeyFunc func(Zettel) string

type Set struct {
	keyFunc  SetKeyFunc
	innerMap map[string]Zettel
}

func MakeSetKeyFuncHinweis() SetKeyFunc {
	return func(sz Zettel) string {
		return makeKey(sz.Named.Hinweis)
	}
}

func makeKey(ss ...fmt.Stringer) string {
	sb := &strings.Builder{}

	for i, s := range ss {
		if i > 0 {
			sb.WriteString(".")
		}

		sb.WriteString(s.String())
	}

	return sb.String()
}
