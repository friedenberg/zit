package sku

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/int_value"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/ts"
)

type Mutter [2]ts.Time

type IdLike = fmt.Stringer

type DataIdentity interface {
	GetTime() ts.Time
	gattung.GattungLike
	GetObjekteSha() sha.Sha
	//TODO-P1 add GetTime
	//TODO-P1 add GetAkteSha
}

type SkuLike interface {
	DataIdentity
	GetId() IdLike
	SetTimeAndFields(ts.Time, ...string) error
	SetFromSku(Sku) error

	GetKey() string

	GetMutter() Mutter
	SetTransactionIndex(int)
	GetTransactionIndex() int_value.IntValue
	GetKopf() ts.Time
	GetSchwanz() ts.Time
}

type Sku struct {
	Gattung gattung.Gattung

	Time ts.Time

	Kennung    collections.StringValue
	ObjekteSha sha.Sha
	AkteSha    sha.Sha
}

func MakeSku(line string) (sk Sku, err error) {
	if err = sk.Set(line); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (sk *Sku) Set(line string) (err error) {
	r := strings.NewReader(line)

	if _, err = format.ReadSep(
		' ',
		r,
		sk.Gattung.Set,
		sk.Time.Set,
		sk.Kennung.Set,
		sk.ObjekteSha.Set,
		sk.AkteSha.Set,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Sku) Reset(b *Sku) {
	if b == nil {
		a.ObjekteSha = sha.Sha{}
		a.AkteSha = sha.Sha{}
	} else {
		a.ObjekteSha = b.ObjekteSha
		a.AkteSha = b.AkteSha
	}
}

func (a Sku) GetTime() ts.Time {
	return a.Time
}

func (a Sku) GetGattung() gattung.Gattung {
	return a.Gattung
}

func (a Sku) GetObjekteSha() sha.Sha {
	return a.ObjekteSha
}

func (a Sku) Less(b *Sku) (ok bool) {
	if a.Time.Less(b.Time) {
		ok = true
		return
	}

	return
}

func (a Sku) Equals(b *Sku) (ok bool) {
	if a != *b {
		return false
	}

	return true
}

func (s Sku) String() string {
	return fmt.Sprintf(
		"%s %s %s %s %s",
		s.Gattung,
		s.Time,
		s.Kennung,
		s.ObjekteSha,
		s.AkteSha,
	)
}
