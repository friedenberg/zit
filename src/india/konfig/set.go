package konfig

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/kasten"
)

func init() {
	gob.RegisterName("typSet", makeCompiledTypSet(nil))
	gob.RegisterName("etikettSet", makeCompiledEtikettSet(nil))
	gob.RegisterName("kastenSet", makeCompiledKastenSet(nil))
}

type set[
	E schnittstellen.Transacted[E],
	EPtr schnittstellen.TransactedPtr[E],
] struct {
	collections.Set2[E, EPtr]
}

type kastenSet = set[kasten.Transacted, *kasten.Transacted]

func makeCompiledKastenSet(s1 collections.SetLike[*kasten.Transacted]) (s kastenSet) {
	s.Set2 = collections.Set2FromSetLike[kasten.Transacted, *kasten.Transacted](s, s1)

	return
}

func makeCompiledKastenSetFromSlice(s1 []*kasten.Transacted) (s kastenSet) {
	s.Set2 = collections.Set2FromSlice[
		kasten.Transacted,
		*kasten.Transacted,
	](s, s1...)

	return
}

type etikettSet = set[etikett.Transacted, *etikett.Transacted]

func makeCompiledEtikettSetFromSlice(s1 []*etikett.Transacted) (s etikettSet) {
	s.Set2 = collections.Set2FromSlice[
		etikett.Transacted,
		*etikett.Transacted,
	](s, s1...)

	return
}

func makeCompiledEtikettSet(s1 collections.SetLike[*etikett.Transacted]) (s etikettSet) {
	s.Set2 = collections.Set2FromSetLike[etikett.Transacted, *etikett.Transacted](s, s1)

	return
}

type typSet = set[typ.Transacted, *typ.Transacted]

func makeCompiledTypSetFromSlice(s1 []*typ.Transacted) (s typSet) {
	s.Set2 = collections.Set2FromSlice[typ.Transacted, *typ.Transacted](s, s1...)

	return
}

func makeCompiledTypSet(s1 collections.SetLike[*typ.Transacted]) (s typSet) {
	s.Set2 = collections.Set2FromSetLike[typ.Transacted, *typ.Transacted](s, s1)

	return
}

func (s set[E, EPtr]) Key(v EPtr) string {
	if v == nil {
		return ""
	}

	return v.GetKennungString()
}
