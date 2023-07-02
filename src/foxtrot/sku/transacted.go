package sku

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

// TODO-P2 move sku.Sku to sku.Transacted
type Transacted[K kennung.KennungLike[K], KPtr kennung.KennungLikePtr[K]] struct {
	WithKennung      metadatei.WithKennung[K, KPtr]
	ObjekteSha       sha.Sha
	TransactionIndex values.Int
	Kopf             kennung.Tai
}

func (t *Transacted[K, KPtr]) SetFromSku(sk Sku) (err error) {
	if err = KPtr(&t.WithKennung.Kennung).Set(sk.WithKennung.Kennung.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.ObjekteSha = sk.ObjekteSha
	t.WithKennung.Metadatei.AkteSha = sk.WithKennung.Metadatei.AkteSha
	t.GetMetadateiPtr().Tai = sk.GetTai()

	t.Kopf = sk.GetTai()

	return
}

func MakeSkuTransacted(t kennung.Tai, line string) (out SkuLikePtr, err error) {
	fields := strings.Fields(line)
	var g gattung.Gattung

	if err = g.Set(fields[0]); err != nil {
		err = errors.Wrapf(err, "failed to set type: %s", fields[0])
		return
	}

	switch g {
	case gattung.Zettel:
		out = &Transacted[kennung.Hinweis, *kennung.Hinweis]{}

	case gattung.Typ:
		out = &Transacted[kennung.Typ, *kennung.Typ]{}

	case gattung.Etikett:
		out = &Transacted[kennung.Etikett, *kennung.Etikett]{}

	case gattung.Konfig:
		out = &Transacted[kennung.Konfig, *kennung.Konfig]{}

	default:
		err = errors.Errorf("unsupported gattung: %s", g)
		return
	}

	if err = out.SetTimeAndFields(t, fields[1:]...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a Transacted[K, KPtr]) String() string {
	return fmt.Sprintf(
		"%s %s %s",
		a.WithKennung.Kennung,
		a.ObjekteSha,
		a.WithKennung.Metadatei.AkteSha,
	)
}

func (a Transacted[K, KPtr]) GetMetadatei() Metadatei {
	return a.WithKennung.GetMetadatei()
}

func (a *Transacted[K, KPtr]) GetMetadateiPtr() *Metadatei {
	return a.WithKennung.GetMetadateiPtr()
}

func (a Transacted[K, KPtr]) GetTai() kennung.Tai {
	return a.GetMetadatei().GetTai()
}

func (a *Transacted[K, KPtr]) SetTai(t kennung.Tai) {
	a.GetMetadateiPtr().Tai = t
}

func (a Transacted[K, KPtr]) GetKennung() K {
	return a.WithKennung.Kennung
}

func (a Transacted[K, KPtr]) GetExternal() External[K, KPtr] {
	return External[K, KPtr]{
		WithKennung: metadatei.WithKennung[K, KPtr]{
			Kennung: a.WithKennung.Kennung,
			Metadatei: Metadatei{
				AkteSha: sha.Make(a.GetAkteSha()),
			},
		},
		ObjekteSha: a.ObjekteSha,
	}
}

func (a *Transacted[K, KPtr]) Sku() Sku {
	return Sku{
		WithKennung: metadatei.WithKennungInterface{
			Kennung: a.WithKennung.Kennung,
			Metadatei: Metadatei{
				Tai:     a.GetTai(),
				Gattung: gattung.Make(a.GetGattung()),
				AkteSha: a.WithKennung.Metadatei.AkteSha,
			},
		},
		ObjekteSha: a.ObjekteSha,
	}
}

func (a *Transacted[K, KPtr]) SetTransactionIndex(i int) {
	a.TransactionIndex.SetInt(i)
}

func (a *Transacted[K, KPtr]) Reset() {
	a.Kopf.Reset()
	a.ObjekteSha.Reset()
	a.WithKennung.Reset()
	a.TransactionIndex.Reset()
}

func (a *Transacted[K, KPtr]) ResetWith(b Transacted[K, KPtr]) {
	a.Kopf = b.Kopf
	a.ObjekteSha = b.ObjekteSha
	a.WithKennung.ResetWith(b.WithKennung)
	a.TransactionIndex.SetInt(b.TransactionIndex.Int())
}

func (a Transacted[K, KPtr]) Less(b Transacted[K, KPtr]) (ok bool) {
	if a.GetTai().Less(b.GetTai()) {
		ok = true
		return
	}

	if a.GetTai().Equals(b.GetTai()) &&
		a.TransactionIndex.Less(b.TransactionIndex) {
		ok = true
		return
	}

	return
}

func (a Transacted[K, KPtr]) EqualsAny(b any) (ok bool) {
	return values.Equals(a, b)
}

func (a Transacted[K, KPtr]) Equals(b Transacted[K, KPtr]) (ok bool) {
	if !a.TransactionIndex.Equals(b.TransactionIndex) {
		return
	}

	if !a.GetTai().Equals(b.GetTai()) {
		return
	}

	if a.GetKennung().String() != b.GetKennung().String() {
		return
	}

	if !a.ObjekteSha.Equals(b.ObjekteSha) {
		return
	}

	if !a.WithKennung.Metadatei.AkteSha.Equals(
		b.WithKennung.Metadatei.AkteSha,
	) {
		return
	}

	return true
}

func (o *Transacted[K, KPtr]) SetTimeAndFields(
	t kennung.Tai,
	vs ...string,
) (err error) {
	o.WithKennung.GetMetadateiPtr().Tai = t

	if len(vs) != 4 {
		err = errors.Errorf("expected 4 elements but got %d", len(vs))
		return
	}

	// Mutter[0] used to be here

	vs = vs[1:]

	// Mutter[1] used to be here

	vs = vs[1:]

	if err = KPtr(&o.WithKennung.Kennung).Set(vs[0]); err != nil {
		err = errors.Wrapf(err, "failed to set id: %s", vs[1])
		return
	}

	vs = vs[1:]

	if err = o.ObjekteSha.Set(vs[0]); err != nil {
		err = errors.Wrapf(err, "failed to set sha: %s", vs[2])
		return
	}

	return
}

func (s Transacted[K, KPtr]) GetGattung() schnittstellen.GattungLike {
	return s.WithKennung.Kennung.GetGattung()
}

func (s Transacted[K, KPtr]) GetId() Kennung {
	return KPtr(&s.WithKennung.Kennung)
}

func (s Transacted[K, KPtr]) GetObjekteSha() schnittstellen.ShaLike {
	return s.ObjekteSha
}

func (s Transacted[K, KPtr]) GetAkteSha() schnittstellen.ShaLike {
	return s.WithKennung.Metadatei.AkteSha
}

func (s Transacted[K, KPtr]) GetTransactionIndex() values.Int {
	return s.TransactionIndex
}

func (o Transacted[K, KPtr]) GetKey() string {
	return fmt.Sprintf("%s.%s", o.GetGattung(), o.GetKennung())
}
