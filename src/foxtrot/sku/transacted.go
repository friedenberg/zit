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
)

// TODO-P2 move sku.Sku to sku.Transacted
type Transacted[T kennung.KennungLike[T], T1 kennung.KennungLikePtr[T]] struct {
	Kennung          T
	ObjekteSha       sha.Sha
	AkteSha          sha.Sha
	TransactionIndex values.Int
	Kopf, Schwanz    kennung.Tai
}

func (t *Transacted[T, T1]) SetFromSku(sk Sku) (err error) {
	t.Schwanz = sk.GetTai()

	if err = T1(&t.Kennung).Set(sk.Kennung.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.ObjekteSha = sk.ObjekteSha
	t.AkteSha = sk.AkteSha

	t.Kopf = sk.GetTai()
	t.Schwanz = sk.GetTai()

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

func (a Transacted[T, T1]) String() string {
	return fmt.Sprintf("%s %s %s", a.Kennung, a.ObjekteSha, a.AkteSha)
}

func (a Transacted[T, T1]) GetTai() kennung.Tai {
	return a.Schwanz
}

func (a Transacted[T, T1]) GetExternal() External[T, T1] {
	return External[T, T1]{
		Kennung:    a.Kennung,
		AkteSha:    a.AkteSha,
		ObjekteSha: a.ObjekteSha,
	}
}

func (a *Transacted[T, T1]) Sku() Sku {
	return Sku{
		Tai:        a.GetTai(),
		Gattung:    gattung.Make(a.GetGattung()),
		Kennung:    values.MakeString(a.Kennung.String()),
		ObjekteSha: a.ObjekteSha,
		AkteSha:    a.AkteSha,
	}
}

func (a *Transacted[T, T1]) SetTransactionIndex(i int) {
	a.TransactionIndex.SetInt(i)
}

func (a *Transacted[T, T1]) Reset() {
	a.Kopf.Reset()
	a.ObjekteSha.Reset()
	a.AkteSha.Reset()
	T1(&a.Kennung).Reset()
	a.Schwanz.Reset()
	a.TransactionIndex.Reset()
}

func (a *Transacted[T, T1]) ResetWith(b Transacted[T, T1]) {
	a.Kopf = b.Kopf
	a.ObjekteSha = b.ObjekteSha
	a.AkteSha = b.AkteSha
	T1(&a.Kennung).ResetWith(b.Kennung)
	a.Schwanz = b.Schwanz
	a.TransactionIndex.SetInt(b.TransactionIndex.Int())
}

func (a Transacted[T, T1]) Less(b Transacted[T, T1]) (ok bool) {
	if a.Schwanz.Less(b.Schwanz) {
		ok = true
		return
	}

	if a.Schwanz.Equals(b.Schwanz) &&
		a.TransactionIndex.Less(b.TransactionIndex) {
		ok = true
		return
	}

	return
}

func (a Transacted[T, T1]) EqualsAny(b any) (ok bool) {
	return values.Equals(a, b)
}

func (a Transacted[T, T1]) Equals(b Transacted[T, T1]) (ok bool) {
	if !a.TransactionIndex.Equals(b.TransactionIndex) {
		return
	}

	if !a.Schwanz.Equals(b.Schwanz) {
		return
	}

	if a.Kennung.String() != b.Kennung.String() {
		return
	}

	if !a.ObjekteSha.Equals(b.ObjekteSha) {
		return
	}

	if !a.AkteSha.Equals(b.AkteSha) {
		return
	}

	return true
}

func (o *Transacted[T, T1]) SetTimeAndFields(
	t kennung.Tai,
	vs ...string,
) (err error) {
	o.Schwanz = t

	if len(vs) != 4 {
		err = errors.Errorf("expected 4 elements but got %d", len(vs))
		return
	}

	// Mutter[0] used to be here

	vs = vs[1:]

	// Mutter[1] used to be here

	vs = vs[1:]

	if err = T1(&o.Kennung).Set(vs[0]); err != nil {
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

func (s Transacted[T, T1]) GetGattung() schnittstellen.Gattung {
	return s.Kennung.GetGattung()
}

func (s Transacted[T, T1]) GetId() Kennung {
	return s.Kennung
}

func (s Transacted[T, T1]) GetObjekteSha() schnittstellen.Sha {
	return s.ObjekteSha
}

func (s Transacted[T, T1]) GetAkteSha() schnittstellen.Sha {
	return s.AkteSha
}

func (s Transacted[T, T1]) GetTransactionIndex() values.Int {
	return s.TransactionIndex
}

func (o Transacted[T, T1]) GetKey() string {
	return fmt.Sprintf("%s.%s", o.Kennung.GetGattung(), o.Kennung)
}
