package konfig

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/hotel/etikett"
)

func init() {
	collections.RegisterGob[ketikett]()
}

type ketikett struct {
	Transacted        etikett.Transacted
	ImplicitEtiketten schnittstellen.MutableSet[kennung.Etikett]
}

func (a ketikett) Less(b ketikett) bool {
	return a.Transacted.Less(b.Transacted)
}

func (a ketikett) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a ketikett) Equals(b ketikett) bool {
	if !a.Transacted.Equals(b.Transacted) {
		return false
	}

	if !a.ImplicitEtiketten.Equals(b.ImplicitEtiketten) {
		return false
	}

	return true
}

func (e ketikett) String() string {
	return e.Transacted.GetKennung().String()
}
