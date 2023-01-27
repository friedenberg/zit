package zettel_external

import (
	"fmt"

	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/fd"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type Sku = sku.External[kennung.Hinweis, *kennung.Hinweis]

// TODO-P3 rename to External?
type Zettel struct {
	Objekte  zettel.Objekte
	Sku      Sku
	ZettelFD fd.FD
	AkteFD   fd.FD
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
