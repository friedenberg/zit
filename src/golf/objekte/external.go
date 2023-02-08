package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
)

type ExternalKeyer[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
	T2 schnittstellen.Id[T2],
	T3 schnittstellen.IdPtr[T2],
] struct{}

func (_ ExternalKeyer[T, T1, T2, T3]) Key(e *External[T, T1, T2, T3]) string {
	if e == nil {
		return ""
	}

	return e.Sku.Kennung.String()
}

type External[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
	T2 schnittstellen.Id[T2],
	T3 schnittstellen.IdPtr[T2],
] struct {
	Objekte T
	Sku     sku.External[T2, T3]
	FD      kennung.FD
}

func (e External[T, T1, T2, T3]) GetGattung() schnittstellen.Gattung {
	return e.Sku.Kennung.GetGattung()
}

func (e External[T, T1, T2, T3]) GetObjekteSha() sha.Sha {
	return e.Sku.ObjekteSha
}

func (e External[T, T1, T2, T3]) GetAkteSha() sha.Sha {
	return e.Sku.AkteSha
}

func (e *External[T, T1, T2, T3]) SetAkteSha(v sha.Sha) {
	e.Sku.ObjekteSha = v
}

func (e External[T, T1, T2, T3]) ObjekteSha() sha.Sha {
	return e.Sku.ObjekteSha
}

func (e *External[T, T1, T2, T3]) SetObjekteSha(
	arf schnittstellen.AkteReaderFactory,
	v string,
) (err error) {
	if err = e.Sku.ObjekteSha.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
