package sku

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/format"
	"github.com/friedenberg/zit/src/foxtrot/sha"
	"github.com/friedenberg/zit/src/golf/ts"
)

type Sku2 struct {
	Gattung gattung.Gattung

	Tai ts.Tai

	Kennung    collections.StringValue
	ObjekteSha sha.Sha
	AkteSha    sha.Sha
}

func MakeSku2(line string) (sk Sku2, err error) {
	if err = sk.Set(line); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (sk *Sku2) Set(line string) (err error) {
	r := strings.NewReader(line)

	if _, err = format.ReadSep(
		' ',
		r,
		sk.Gattung.Set,
		sk.Tai.Set,
		sk.Kennung.Set,
		sk.ObjekteSha.Set,
		sk.AkteSha.Set,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Sku2) Reset(b *Sku2) {
	if b == nil {
		a.ObjekteSha = sha.Sha{}
		a.AkteSha = sha.Sha{}
	} else {
		a.ObjekteSha = b.ObjekteSha
		a.AkteSha = b.AkteSha
	}
}

func (a Sku2) GetTai() ts.Tai {
	return a.Tai
}

func (a Sku2) GetKopf() ts.Time {
	return a.Tai.AsTime()
}

func (a Sku2) GetSchwanz() ts.Time {
	return a.Tai.AsTime()
}

func (a Sku2) GetKey() string {
	return a.String()
}

func (a Sku2) GetTime() ts.Time {
	return a.Tai.AsTime()
}

func (a Sku2) GetId() IdLike {
	return a.Kennung
}

func (a Sku2) GetGattung() gattung.Gattung {
	return a.Gattung
}

func (a Sku2) GetObjekteSha() sha.Sha {
	return a.ObjekteSha
}

func (a Sku2) GetAkteSha() sha.Sha {
	return a.AkteSha
}

func (a Sku2) Less(b Sku2) (ok bool) {
	if a.Tai.Less(b.Tai) {
		ok = true
		return
	}

	return
}

func (a Sku2) Equals(b *Sku2) (ok bool) {
	if b == nil {
		return false
	}

	if a != *b {
		return false
	}

	return true
}

func (s Sku2) String() string {
	return fmt.Sprintf(
		"%s %s %s %s %s",
		s.Gattung,
		s.Tai,
		s.Kennung,
		s.ObjekteSha,
		s.AkteSha,
	)
}
