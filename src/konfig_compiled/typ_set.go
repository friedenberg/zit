package konfig_compiled

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/golf/typ"
)

func init() {
	gob.Register(makeCompiledTypSet(nil))
}

type typSet struct {
	collections.Set2[typ.Transacted, *typ.Transacted]
}

func makeCompiledTypSetFromSlice(s1 []*typ.Transacted) (s typSet) {
	s.Set2 = collections.Set2FromSlice[typ.Transacted, *typ.Transacted](s, s1...)

	return
}

func makeCompiledTypSet(s1 collections.SetLike[*typ.Transacted]) (s typSet) {
	s.Set2 = collections.Set2FromSetLike[typ.Transacted, *typ.Transacted](s, s1)

	return
}

func (s typSet) Key(v *typ.Transacted) string {
	if v == nil {
		return ""
	}

	return v.Sku.Kennung.String()
}
