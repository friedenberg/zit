package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/values"
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
}

func (a External[T, T1, T2, T3]) GetEtiketten() kennung.EtikettSet {
	egs := []any{
		a.Objekte,
	}

	for _, o := range egs {
		if eg, ok := o.(kennung.EtikettenGetter); ok {
			return eg.GetEtiketten()
		}
	}

	return kennung.MakeEtikettSet()
}

func (a External[T, T1, T2, T3]) GetEtikettenExpanded() kennung.EtikettSet {
	egs := []any{
		a.Objekte,
	}

	for _, o := range egs {
		if eg, ok := o.(kennung.EtikettenExpandedGetter); ok {
			return eg.GetEtikettenExpanded()
		}
	}

	return kennung.Expanded(a.GetEtiketten())
}

func (a External[T, T1, T2, T3]) GetTyp() kennung.Typ {
	ok := false
	o := any(a.Objekte)

	var tg kennung.TypGetter

	if tg, ok = o.(kennung.TypGetter); !ok {
		tg = nil
		return kennung.MustTyp(a.GetGattung().GetGattungString())
	}

	return tg.GetTyp()
}

func (a External[T, T1, T2, T3]) GetIdLike() (il kennung.IdLike) {
	return a.Sku.Kennung
}

func (a External[T, T1, T2, T3]) String() string {
	return a.Sku.String()
}

func (a External[T, T1, T2, T3]) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a External[T, T1, T2, T3]) Equals(b External[T, T1, T2, T3]) bool {
	if !a.Objekte.Equals(b.Objekte) {
		return false
	}

	if !a.Sku.Equals(b.Sku) {
		return false
	}

	return true
}

func (e External[T, T1, T2, T3]) GetGattung() schnittstellen.Gattung {
	return e.Sku.Kennung.GetGattung()
}

func (e External[T, T1, T2, T3]) GetObjekteFD() kennung.FD {
	return e.Sku.FDs.Objekte
}

func (e External[T, T1, T2, T3]) GetAkteFD() kennung.FD {
	return e.Sku.FDs.Akte
}

func (e External[T, T1, T2, T3]) GetObjekteSha() schnittstellen.Sha {
	return e.Sku.ObjekteSha
}

func (e External[T, T1, T2, T3]) GetAkteSha() schnittstellen.Sha {
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
