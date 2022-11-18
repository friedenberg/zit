package string_expansion

import (
	"regexp"

	"github.com/friedenberg/zit/src/bravo/collections"
)

type expanderRight[T collections.ProtoObjekte, T1 interface {
	*T
	collections.ProtoObjektePointer
}] struct {
	delimiter *regexp.Regexp
}

func MakeExpanderRight[T collections.ProtoObjekte, T1 interface {
	*T
	collections.ProtoObjektePointer
}](delimiter string) expanderRight[T, T1] {
	return expanderRight[T, T1]{
		delimiter: regexp.MustCompile(delimiter),
	}
}

func (ex expanderRight[T, T1]) Expand(s string) (out collections.ValueSet[T, T1]) {
	expanded := collections.MakeMutableValueSet[T, T1]()
	expanded.AddString(s)

	defer func() {
		out = expanded.Copy()
	}()

	if s == "" {
		return
	}

	hyphens := ex.delimiter.FindAllIndex([]byte(s), -1)

	if hyphens == nil {
		return
	}

	for _, loc := range hyphens {
		locStart := loc[0]
		t1 := s[0:locStart]

		expanded.AddString(t1)
	}

	return
}
