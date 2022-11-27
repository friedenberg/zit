package sku

import (
	"flag"
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/int_value"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/delta/ts"
)

type Sku2 struct {
	Gattung          gattung.Gattung
	Mutter           Mutter
	Id               flag.Value
	Sha              sha.Sha
	TransactionIndex int_value.IntValue
	Kopf, Schwanz    ts.Time
}

func (a *Sku2) Reset(b *Sku2) {
	if b == nil {
		a.Kopf = b.Kopf
		a.Mutter[0] = b.Mutter[0]
		a.Mutter[1] = b.Mutter[1]
		a.Schwanz = b.Schwanz
		a.TransactionIndex.SetInt(b.TransactionIndex.Int())
	} else {
		a.Kopf = ts.Time{}
		a.Mutter[0] = ts.Time{}
		a.Mutter[1] = ts.Time{}
		a.Schwanz = ts.Time{}
		a.TransactionIndex.Reset()
	}
}

func (a Sku2) Less(b *Sku2) (ok bool) {
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

func (a Sku2) Equals(b *Sku2) (ok bool) {
	if !a.TransactionIndex.Equals(b.TransactionIndex) {
		return
	}

	if !a.Schwanz.Equals(b.Schwanz) {
		return
	}

	if a.Gattung != b.Gattung {
		return
	}

	if a.Mutter != b.Mutter {
		return
	}

	if a.Id.String() != b.Id.String() {
		return
	}

	if !a.Sha.Equals(b.Sha) {
		return
	}

	return true
}

func (o *Sku2) Set(v string) (err error) {
	vs := strings.Split(v, " ")

	if len(vs) != 5 {
		err = errors.Errorf("expected 5 elements but got %d", len(vs))
		return
	}

	if err = o.Gattung.Set(vs[0]); err != nil {
		err = errors.Wrapf(err, "failed to set type: %s", vs[0])
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

	switch o.Gattung {
	case gattung.Zettel:
		o.Id = &hinweis.Hinweis{}

	case gattung.Etikett:
		o.Id = &kennung.Etikett{}

	case gattung.Typ:
		o.Id = &kennung.Typ{}

	case gattung.Konfig:
		o.Id = &kennung.Konfig{}

	default:
		err = errors.Errorf("unsupported gattung: %s", o.Gattung)
		return
	}

	if err = o.Id.Set(vs[0]); err != nil {
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

func (o Sku2) GetKey() string {
	return fmt.Sprintf("%s.%s", o.Gattung, o.Id)
}
