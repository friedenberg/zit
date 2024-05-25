package metadatei

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

var (
	EqualerSansTai               equaler
	Equaler                      = equaler{includeTai: true}
	EqualerSansTaiIncludeVirtual = equaler{includeVirtual: true}
)

type equaler struct {
	includeVirtual bool
	includeTai     bool
}

func (e equaler) Equals(a, b *Metadatei) bool {
	if e.includeTai && !a.Tai.Equals(b.Tai) {
		// log.Debug().Print(&a.Tai, "->", &b.Tai)
		return false
	}

	if !a.Akte.EqualsSha(&b.Akte) {
		// log.Debug().Print(&a.Akte, "->", &b.Akte)
		return false
	}

	if !a.Typ.Equals(b.Typ) {
		// log.Debug().Print(&a.Typ, "->", &b.Typ)
		return false
	}

	aes := a.GetEtiketten()
	bes := b.GetEtiketten()

  // for i, ea := range a.Verzeichnisse.Etiketten.All {

  // }

	if err := aes.EachPtr(
		func(ea *kennung.Etikett) (err error) {
			if !e.includeVirtual && ea.IsVirtual() {
				return
			}

			if !bes.ContainsKey(bes.KeyPtr(ea)) {
				err = errors.New("false")
				return
			}

			return
		},
	); err != nil {
		return false
	}

	if err := bes.EachPtr(
		func(eb *kennung.Etikett) (err error) {
			if !e.includeVirtual && eb.IsVirtual() {
				return
			}

			if !aes.ContainsKey(aes.KeyPtr(eb)) {
				err = errors.New("false")
				return
			}

			return
		},
	); err != nil {
		// log.Debug().Print(aes, "->", bes)
		return false
	}

	if !a.Bezeichnung.Equals(b.Bezeichnung) {
		// log.Debug().Print(a.Bezeichnung, "->", b.Bezeichnung)
		return false
	}

	return true
}
