package sku

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/int_value"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/ts"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
)

type SkuLikeOld interface {
	SkuLike
	SetTimeAndFields(ts.Time, ...string) error
}

// TODO-P3 examine adding objekte and akte shas to Skus
// TODO-P2 move sku.Sku to sku.Transacted
type Transacted[T kennung.KennungLike[T], T1 kennung.KennungLikePtr[T]] struct {
	Kennung    T
	ObjekteSha sha.Sha
	AkteSha    sha.Sha
	Verzeichnisse
}

func (t *Transacted[T, T1]) SetFromSku2(sk Sku2) (err error) {
	t.Schwanz = sk.GetTime()

	if err = T1(&t.Kennung).Set(sk.Kennung.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.ObjekteSha = sk.ObjekteSha
	t.AkteSha = sk.AkteSha

	t.Verzeichnisse.SetFromSku2(sk)

	return
}

func (t *Transacted[T, T1]) SetFromSku(sk Sku) (err error) {
	t.Schwanz = sk.Time

	if err = T1(&t.Kennung).Set(sk.Kennung.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.ObjekteSha = sk.ObjekteSha
	t.AkteSha = sk.AkteSha

	t.Verzeichnisse.SetFromSku(sk)

	return
}

// func TransactedFromSku(sk Sku) (out SkuLike, err error) {
// 	switch sk.Gattung {
// 	case gattung.Zettel:
// 		out = &Transacted[hinweis.Hinweis, *hinweis.Hinweis]{}

// 	case gattung.Typ:
// 		out = &Transacted[kennung.Typ, *kennung.Typ]{}

// 	case gattung.Etikett:
// 		out = &Transacted[kennung.Etikett, *kennung.Etikett]{}

// 	case gattung.Konfig:
// 		out = &Transacted[kennung.Konfig, *kennung.Konfig]{}

// 	default:
// 		err = errors.Errorf("unsupported gattung: %s", sk.Gattung)
// 		return
// 	}

// 	if err = out.SetFromSku(sk); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

// TODO-P2 include sku versions
func MakeSkuTransacted(t ts.Time, line string) (out SkuLike, err error) {
	fields := strings.Fields(line)
	var g gattung.Gattung

	if err = g.Set(fields[0]); err != nil {
		err = errors.Wrapf(err, "failed to set type: %s", fields[0])
		return
	}

	switch g {
	case gattung.Zettel:
		out = &Transacted[hinweis.Hinweis, *hinweis.Hinweis]{}

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

func (a Transacted[T, T1]) GetTime() ts.Time {
	return a.Schwanz
}

func (a *Transacted[T, T1]) Sku() Sku {
	return Sku{
		Time:       ts.TimeWithIndex(a.GetTime(), a.GetTransactionIndex().Int()),
		Gattung:    gattung.Make(a.GetGattung()),
		Kennung:    collections.MakeStringValue(a.Kennung.String()),
		ObjekteSha: a.ObjekteSha,
		AkteSha:    a.AkteSha,
	}
}

func (a *Transacted[T, T1]) Sku2() Sku2 {
	return Sku2{
		Tai: ts.TaiFromTimeWithIndex(
			a.GetTime(),
			a.GetTransactionIndex().Int(),
		),
		Gattung:    gattung.Make(a.GetGattung()),
		Kennung:    collections.MakeStringValue(a.Kennung.String()),
		ObjekteSha: a.ObjekteSha,
		AkteSha:    a.AkteSha,
	}
}

func (a *Transacted[T, T1]) SetTransactionIndex(i int) {
	a.TransactionIndex.SetInt(i)
}

func (a *Transacted[T, T1]) Reset() {
	a.Kopf = ts.Time{}
	a.ObjekteSha = sha.Sha{}
	a.AkteSha = sha.Sha{}
	T1(&a.Kennung).Reset()
	a.Mutter[0] = ts.Time{}
	a.Mutter[1] = ts.Time{}
	a.Schwanz = ts.Time{}
	a.TransactionIndex.Reset()
}

func (a *Transacted[T, T1]) ResetWith(b Transacted[T, T1]) {
	a.Kopf = b.Kopf
	a.ObjekteSha = b.ObjekteSha
	a.AkteSha = b.AkteSha
	T1(&a.Kennung).ResetWith(b.Kennung)
	a.Mutter[0] = b.Mutter[0]
	a.Mutter[1] = b.Mutter[1]
	a.Schwanz = b.Schwanz
	a.TransactionIndex.SetInt(b.TransactionIndex.Int())
}

func (a Transacted[T, T1]) Less(b *Transacted[T, T1]) (ok bool) {
	if a.Schwanz.Less(b.Schwanz) {
		ok = true
		return
	}

	if a.Schwanz.Equals(b.Schwanz) && a.TransactionIndex.Less(b.TransactionIndex) {
		ok = true
		return
	}

	return
}

func (a Transacted[T, T1]) Equals(b *Transacted[T, T1]) (ok bool) {
	if !a.TransactionIndex.Equals(&b.TransactionIndex) {
		return
	}

	if !a.Schwanz.Equals(b.Schwanz) {
		return
	}

	if a.Mutter != b.Mutter {
		return
	}

	if a.Kennung.String() != b.Kennung.String() {
		return
	}

	if !a.ObjekteSha.Equals(b.ObjekteSha) {
		return
	}

	return true
}

func (o *Transacted[T, T1]) SetTimeAndFields(t ts.Time, vs ...string) (err error) {
	o.Schwanz = t

	if len(vs) != 4 {
		err = errors.Errorf("expected 4 elements but got %d", len(vs))
		return
	}

	if err = o.Mutter[0].Set(vs[0]); err != nil {
		err = errors.Wrapf(err, "failed to set mutter 0: %s", vs[0])
		return
	}

	vs = vs[1:]

	if err = o.Mutter[1].Set(vs[0]); err != nil {
		err = errors.Wrapf(err, "failed to set mutter 1: %s", vs[0])
		return
	}

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

func (s Transacted[T, T1]) GetMutter() Mutter {
	return s.Mutter
}

func (s Transacted[T, T1]) GetGattung() schnittstellen.Gattung {
	return s.Kennung.GetGattung()
}

func (s Transacted[T, T1]) GetId() IdLike {
	return s.Kennung
}

func (s Transacted[T, T1]) GetObjekteSha() schnittstellen.Sha {
	return s.ObjekteSha
}

func (s Transacted[T, T1]) GetAkteSha() schnittstellen.Sha {
	return s.AkteSha
}

func (s Transacted[T, T1]) GetTransactionIndex() int_value.IntValue {
	return s.TransactionIndex
}

func (s Transacted[T, T1]) GetKopf() ts.Time {
	return s.Kopf
}

func (s Transacted[T, T1]) GetSchwanz() ts.Time {
	return s.Schwanz
}

func (o Transacted[T, T1]) GetKey() string {
	return fmt.Sprintf("%s.%s", o.Kennung.GetGattung(), o.Kennung)
}
