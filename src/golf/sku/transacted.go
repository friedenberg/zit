package sku

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/int_value"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/foxtrot/ts"
)

// TODO-P3 examine adding objekte and akte shas to Skus
// TODO-P2 move sku.Sku to sku.Transacted
type Transacted[T kennung.KennungLike[T], T1 kennung.KennungLikePtr[T]] struct {
	Mutter     Mutter
	Kennung    T
	ObjekteSha sha.Sha
	AkteSha    sha.Sha
	//TODO-P2 add verzeichnisse
	TransactionIndex int_value.IntValue
	Kopf, Schwanz    ts.Time
}

// TODO-P2 include sku versions
func MakeSku(line string) (out SkuLike, err error) {
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

	if err = out.SetFields(fields[1:]...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Transacted[T, T1]) SetTransactionIndex(i int) {
	a.TransactionIndex.SetInt(i)
}

func (a *Transacted[T, T1]) Reset(b *Transacted[T, T1]) {
	if b == nil {
		a.Kopf = ts.Time{}
		a.ObjekteSha = sha.Sha{}
		a.AkteSha = sha.Sha{}
		T1(&a.Kennung).Reset(nil)
		a.Mutter[0] = ts.Time{}
		a.Mutter[1] = ts.Time{}
		a.Schwanz = ts.Time{}
		a.TransactionIndex.Reset()
	} else {
		a.Kopf = b.Kopf
		a.ObjekteSha = b.ObjekteSha
		a.AkteSha = b.AkteSha
		T1(&a.Kennung).Reset(&b.Kennung)
		a.Mutter[0] = b.Mutter[0]
		a.Mutter[1] = b.Mutter[1]
		a.Schwanz = b.Schwanz
		a.TransactionIndex.SetInt(b.TransactionIndex.Int())
	}
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
	if !a.TransactionIndex.Equals(b.TransactionIndex) {
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

func (o *Transacted[T, T1]) SetFields(vs ...string) (err error) {
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

func (s Transacted[T, T1]) GetGattung() gattung.Gattung {
	return s.Kennung.Gattung()
}

func (s Transacted[T, T1]) GetId() IdLike {
	return s.Kennung
}

func (s Transacted[T, T1]) GetObjekteSha() sha.Sha {
	return s.ObjekteSha
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
	return fmt.Sprintf("%s.%s", o.Kennung.Gattung(), o.Kennung)
}
