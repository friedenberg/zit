package metadatei

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

var (
	EqualerSansTai equalerSansTai
	Equaler        equaler
)

type equalerSansTai struct{}

func (equalerSansTai) Equals(a, b *Metadatei) bool {
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

	if err := aes.EachPtr(
		func(ea *kennung.Etikett) (err error) {
			if ea.IsVirtual() {
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
			if eb.IsVirtual() {
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

type equaler struct{}

func (equaler) Equals(pz, z1 *Metadatei) bool {
	if !EqualerSansTai.Equals(pz, z1) {
		return false
	}

	if !pz.Tai.Equals(z1.Tai) {
		return false
	}

	return true
}
