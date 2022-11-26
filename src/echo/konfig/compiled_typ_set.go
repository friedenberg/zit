package konfig

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/bravo/collections"
)

func init() {
	gob.Register(makeCompiledTypSet(nil))
}

type compiledTypSet struct {
	collections.Set2[compiledTyp, *compiledTyp]
}

func makeCompiledTypSet(s1 collections.SetLike[*compiledTyp]) (s compiledTypSet) {
	s.Set2 = collections.Set2FromSetLike[compiledTyp, *compiledTyp](s, s1)

	return
}

func (s compiledTypSet) Key(v *compiledTyp) string {
	if v == nil {
		return ""
	}

	return v.Name.String()
}
