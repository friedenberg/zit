package zettel_external

import (
	"fmt"

	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type Sku = sku.External[kennung.Hinweis, *kennung.Hinweis]

// TODO-P3 rename to External?
type Zettel struct {
	Objekte  zettel.Objekte
	Sku      Sku
	ZettelFD kennung.FD
	AkteFD   kennung.FD
}

func (a Zettel) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Zettel) Equals(b Zettel) bool {
	if !a.Objekte.Equals(b.Objekte) {
		return false
	}

	if !a.Sku.Equals(b.Sku) {
		return false
	}

	if !a.ZettelFD.Equals(b.ZettelFD) {
		return false
	}

	if !a.AkteFD.Equals(b.AkteFD) {
		return false
	}

	return true
}

func (e Zettel) GetObjekteFD() kennung.FD {
	return e.ZettelFD
}

func (e Zettel) GetAkteFD() kennung.FD {
	return e.AkteFD
}

func (e Zettel) GetObjekteSha() sha.Sha {
	return e.Sku.ObjekteSha
}

func (e Zettel) String() string {
	return e.ExternalPathAndSha()
}

func (e Zettel) ExternalPathAndSha() string {
	if !e.ZettelFD.IsEmpty() {
		return fmt.Sprintf("[%s %s]", e.ZettelFD.Path, e.Sku.ObjekteSha)
	} else if !e.AkteFD.IsEmpty() {
		return fmt.Sprintf("[%s %s]", e.AkteFD.Path, e.Objekte.Akte)
	} else {
		return ""
	}
}
