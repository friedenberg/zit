package etikett

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/fd"
	"github.com/friedenberg/zit/src/foxtrot/sku"
)

type ExternalKeyer struct{}

func (_ ExternalKeyer) Key(e *External) string {
	if e == nil {
		return ""
	}

	return e.Sku.Kennung.String()
}

type External struct {
	Objekte Objekte
	Sku     sku.External[kennung.Etikett, *kennung.Etikett]
	FD      fd.FD
}

func (e External) GetGattung() gattung.Gattung {
	return gattung.Etikett
}

func (e External) GetObjekteSha() sha.Sha {
	return e.Sku.ObjekteSha
}

func (e External) GetAkteSha() sha.Sha {
	return e.Objekte.Sha
}

func (e *External) SetAkteSha(v sha.Sha) {
	e.Objekte.Sha = v
}

func (e External) ObjekteSha() sha.Sha {
	return e.Objekte.Sha
}

func (e *External) SetObjekteSha(
	arf schnittstellen.AkteReaderFactory,
	v string,
) (err error) {
	if err = e.Sku.ObjekteSha.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
