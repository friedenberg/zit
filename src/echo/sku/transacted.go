package sku

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/int_value"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/delta/sha"
	"github.com/friedenberg/zit/src/delta/ts"
)

// TODO-P3 examine adding objekte and akte shas to Skus
type Transacted[T kennung.KennungLike[T], T1 kennung.KennungLikePtr[T]] struct {
	Mutter           Mutter
	Kennung          T
	Sha              sha.Sha
	TransactionIndex int_value.IntValue
	Kopf, Schwanz    ts.Time
}

func (a Transacted[T, T1]) Sku() Sku {
	return Sku{
		Gattung: a.Kennung.Gattung(),
		Mutter:  a.Mutter,
		Id:      T1(&a.Kennung),
		Sha:     a.Sha,
	}
}

func (a *Transacted[T, T1]) SetTransactionIndex(i int) {
	a.TransactionIndex.SetInt(i)
}

func (a *Transacted[T, T1]) Reset(b *Transacted[T, T1]) {
	if b == nil {
		a.Kopf = ts.Time{}
		a.Sha = sha.Sha{}
		T1(&a.Kennung).Reset(nil)
		a.Mutter[0] = ts.Time{}
		a.Mutter[1] = ts.Time{}
		a.Schwanz = ts.Time{}
		a.TransactionIndex.Reset()
	} else {
		a.Kopf = b.Kopf
		a.Sha = b.Sha
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

	if !a.Sha.Equals(b.Sha) {
		return
	}

	return true
}

func (o *Transacted[T, T1]) Set(v string) (err error) {
	vs := strings.Split(v, " ")

	if len(vs) != 5 {
		err = errors.Errorf("expected 5 elements but got %d", len(vs))
		return
	}

	var g gattung.Gattung

	if err = g.Set(vs[0]); err != nil {
		err = errors.Wrapf(err, "failed to set type: %s", vs[0])
		return
	}

	if g != o.Kennung.Gattung() {
		err = errors.Errorf(
			"expected gattung %s but got %s",
			o.Kennung.Gattung(),
			g,
		)

		return
	}

	vs = vs[1:]

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

	if err = o.Sha.Set(vs[0]); err != nil {
		err = errors.Wrapf(err, "failed to set sha: %s", vs[2])
		return
	}

	return
}

func (o Transacted[T, T1]) GetKey() string {
	return fmt.Sprintf("%s.%s", o.Kennung.Gattung(), o.Kennung)
}
