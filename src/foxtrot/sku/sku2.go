package sku

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/echo/ts"
)

type Sku2 struct {
	Gattung gattung.Gattung

	Tai ts.Tai

	Kennung    values.String
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
		format.MakeLineReaderIterateStrict(
			sk.Tai.Set,
			sk.Gattung.Set,
			sk.Kennung.Set,
			sk.ObjekteSha.Set,
			sk.AkteSha.Set,
		),
	); err != nil {
		if gattung.IsErrUnrecognizedGattung(err) {
			err = sk.setOld(line)
		} else {
			err = errors.Wrapf(err, "Sku2: %s", line)
		}

		return
	}

	return
}

func (sk *Sku2) setOld(line string) (err error) {
	r := strings.NewReader(line)

	if _, err = format.ReadSep(
		' ',
		r,
		format.MakeLineReaderIterateStrict(
			sk.Gattung.Set,
			sk.Tai.Set,
			sk.Kennung.Set,
			sk.ObjekteSha.Set,
			sk.AkteSha.Set,
		),
	); err != nil {
		err = errors.Wrapf(err, "Sku2: %s", line)
		return
	}

	return
}

func (a *Sku2) ResetWith(b Sku2) {
	errors.TodoP4("should these be more ResetWith calls?")
	a.Gattung = b.Gattung
	a.Tai = b.Tai
	a.Kennung = b.Kennung
	a.ObjekteSha = b.ObjekteSha
	a.AkteSha = b.AkteSha
}

func (a *Sku2) Reset() {
	a.ObjekteSha = sha.Sha{}
	a.AkteSha = sha.Sha{}
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

func (a Sku2) GetGattung() schnittstellen.Gattung {
	return a.Gattung
}

func (a Sku2) GetObjekteSha() schnittstellen.Sha {
	return a.ObjekteSha
}

func (a Sku2) GetAkteSha() schnittstellen.Sha {
	return a.AkteSha
}

func (a Sku2) Less(b Sku2) (ok bool) {
	if a.Tai.Less(b.Tai) {
		ok = true
		return
	}

	return
}

func (a Sku2) EqualsAny(b any) (ok bool) {
	return values.Equals(a, b)
}

func (a Sku2) Equals(b Sku2) (ok bool) {
	if a != b {
		return false
	}

	return true
}

func (s Sku2) String() string {
	return fmt.Sprintf(
		"%s %s %s %s %s",
		s.Tai,
		s.Gattung,
		s.Kennung,
		s.ObjekteSha,
		s.AkteSha,
	)
}
