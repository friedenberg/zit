package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/foxtrot/sku"
)

type ExternalKeyer[
	T Objekte[T],
	T1 ObjektePtr[T],
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
	T Objekte[T],
	T1 ObjektePtr[T],
	T2 schnittstellen.Id[T2],
	T3 schnittstellen.IdPtr[T2],
] struct {
	Objekte   T
	Metadatei metadatei.Metadatei
	Sku       sku.External[T2, T3]
}

func (a External[T, T1, T2, T3]) GetMetadatei() metadatei.Metadatei {
	if mg, ok := any(a.Objekte).(metadatei.Getter); ok {
		return mg.GetMetadatei()
	}

	return a.Metadatei
}

func (a *External[T, T1, T2, T3]) SetMetadatei(m metadatei.Metadatei) {
	if ms, ok := any(&a.Objekte).(metadatei.Setter); ok {
		ms.SetMetadatei(m)
		return
	}

	a.Metadatei = m
}

func (a External[T, T1, T2, T3]) GetEtiketten() kennung.EtikettSet {
	egs := []any{
		// a.Verzeichnisse,
		a.Objekte,
		a.GetMetadatei(),
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
		a.GetMetadatei(),
	}

	for _, o := range egs {
		if eg, ok := o.(kennung.EtikettenExpandedGetter); ok {
			return eg.GetEtikettenExpanded()
		}
	}

	return kennung.Expanded(a.GetEtiketten())
}

func (a External[T, T1, T2, T3]) GetTyp() (t kennung.Typ) {
	tgs := []any{
		// a.Verzeichnisse,
		a.Objekte,
		a.GetMetadatei(),
	}

	for _, o := range tgs {
		if tg, ok := o.(kennung.TypGetter); ok {
			t = tg.GetTyp()
			return
		}
	}

	return
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

func (e External[T, T1, T2, T3]) GetAktePath() string {
	return e.Sku.FDs.Akte.Path
}

func (e External[T, T1, T2, T3]) GetObjekteSha() schnittstellen.Sha {
	return e.Sku.ObjekteSha
}

func (t External[T, T1, T2, T3]) GetAkteSha() schnittstellen.Sha {
	shSku := t.Sku.AkteSha
	shMetadatei := t.GetMetadatei().AkteSha

	if !shSku.Equals(shMetadatei) {
		panic(errors.Errorf(
			"akte sha in sku was %s while akte sha in metadatei was %s",
			shSku,
			shMetadatei,
		))
	}

	return shSku
}

func (e *External[T, T1, T2, T3]) SetAkteSha(v schnittstellen.Sha) {
	sh := sha.Make(v)
	m := e.GetMetadatei()
	m.AkteSha = sh
	e.SetMetadatei(m)
	e.Sku.AkteSha = sh
}

func (e External[T, T1, T2, T3]) ObjekteSha() sha.Sha {
	return e.Sku.ObjekteSha
}

func (e *External[T, T1, T2, T3]) SetObjekteSha(
	sh schnittstellen.Sha,
) {
	e.Sku.ObjekteSha = sha.Make(sh)
}
