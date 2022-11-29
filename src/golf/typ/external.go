package typ

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/delta/typ_toml"
	"github.com/friedenberg/zit/src/echo/fd"
)

type ExternalKeyer struct{}

func (_ ExternalKeyer) Key(e *External) string {
	if e == nil {
		return ""
	}

	return e.Kennung.String()
}

type External struct {
	Objekte typ_toml.Objekte
	Kennung kennung.Typ
	Sha     sha.Sha
	FD      fd.FD
}

func (e External) Gattung() gattung.Gattung {
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
	if err = e.Sha.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
