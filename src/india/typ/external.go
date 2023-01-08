package typ

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/golf/fd"
	"github.com/friedenberg/zit/src/golf/sku"
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
	Sku     sku.External[kennung.Typ, *kennung.Typ]
	FD      fd.FD
}

func (e External) GetGattung() gattung.Gattung {
	return gattung.Typ
}

func (e External) AkteSha() sha.Sha {
	return e.Objekte.Sha
}

func (e *External) SetAkteSha(v sha.Sha) {
	e.Objekte.Sha = v
}

func (e External) ObjekteSha() sha.Sha {
	return e.Objekte.Sha
}

func (e *External) SetObjekteSha(
	arf gattung.AkteReaderFactory,
	v string,
) (err error) {
	if err = e.Sku.ObjekteSha.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
