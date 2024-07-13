package metadatei

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
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

const debug = false

// TODO make better diffing facility
func (e equaler) Equals(a, b *Metadatei) bool {
	if e.includeTai && !a.Tai.Equals(b.Tai) {
		if debug {
			ui.Debug().Print(&a.Tai, "->", &b.Tai)
		}
		return false
	}

	if !a.Akte.EqualsSha(&b.Akte) {
		if debug {
			ui.Debug().Print(&a.Akte, "->", &b.Akte)
		}
		return false
	}

	if !a.Typ.Equals(b.Typ) {
		if debug {
			ui.Debug().Print(&a.Typ, "->", &b.Typ)
		}
		return false
	}

	aes := a.GetEtiketten()
	bes := b.GetEtiketten()

	// for i, ea := range a.Verzeichnisse.Etiketten.All {

	// }

	if err := aes.EachPtr(
		func(ea *kennung.Tag) (err error) {
			if (!e.includeVirtual && ea.IsVirtual()) || ea.IsEmpty() {
				return
			}

			if !bes.ContainsKey(bes.KeyPtr(ea)) {
				if debug {
					ui.Debug().Print(ea, "-> X")
				}

				err = errors.New("false")
				return
			}

			return
		},
	); err != nil {
		if debug {
			ui.Debug().Print(aes, "->", bes)
		}

		return false
	}

	if err := bes.EachPtr(
		func(eb *kennung.Tag) (err error) {
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
		if debug {
			ui.Debug().Print(aes, "->", bes)
		}
		return false
	}

	if !a.Bezeichnung.Equals(b.Bezeichnung) {
		if debug {
			ui.Debug().Print(a.Bezeichnung, "->", b.Bezeichnung)
		}
		return false
	}

	return true
}
