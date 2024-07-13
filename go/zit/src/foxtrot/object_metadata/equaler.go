package object_metadata

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
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
func (e equaler) Equals(a, b *Metadata) bool {
	if e.includeTai && !a.Tai.Equals(b.Tai) {
		if debug {
			ui.Debug().Print(&a.Tai, "->", &b.Tai)
		}
		return false
	}

	if !a.Blob.EqualsSha(&b.Blob) {
		if debug {
			ui.Debug().Print(&a.Blob, "->", &b.Blob)
		}
		return false
	}

	if !a.Type.Equals(b.Type) {
		if debug {
			ui.Debug().Print(&a.Type, "->", &b.Type)
		}
		return false
	}

	aes := a.GetEtiketten()
	bes := b.GetEtiketten()

	// for i, ea := range a.Verzeichnisse.Etiketten.All {

	// }

	if err := aes.EachPtr(
		func(ea *ids.Tag) (err error) {
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
		func(eb *ids.Tag) (err error) {
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

	if !a.Description.Equals(b.Description) {
		if debug {
			ui.Debug().Print(a.Description, "->", b.Description)
		}
		return false
	}

	return true
}
