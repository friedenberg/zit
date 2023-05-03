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
	"github.com/friedenberg/zit/src/delta/kennung"
)

type Sku struct {
	Gattung    gattung.Gattung
	Tai        kennung.Tai
	Kennung    values.String
	ObjekteSha sha.Sha
	AkteSha    sha.Sha
}

func (sk *Sku) Set(line string) (err error) {
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
		if err1 := sk.setOld(line); err1 != nil {
			err = errors.MakeMulti(err, err1)
			return
		}

		err = nil

		return
	}

	return
}

func (sk *Sku) setOld(line string) (err error) {
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

func (a *Sku) ResetWith(b Sku) {
	errors.TodoP4("should these be more ResetWith calls?")
	a.Gattung = b.Gattung
	a.Tai = b.Tai
	a.Kennung = b.Kennung
	a.ObjekteSha = b.ObjekteSha
	a.AkteSha = b.AkteSha
}

func (a *Sku) Reset() {
	a.Gattung.Reset()
	a.Tai.Reset()
	a.Kennung.Reset()
	a.ObjekteSha.Reset()
	a.AkteSha.Reset()
}

func (a Sku) GetTai() kennung.Tai {
	return a.Tai
}

func (a Sku) GetKey() string {
	return a.String()
}

func (a Sku) GetTime() kennung.Time {
	return a.Tai.AsTime()
}

func (a Sku) GetId() Kennung {
	return a.Kennung
}

func (a Sku) GetGattung() schnittstellen.Gattung {
	return a.Gattung
}

func (a Sku) GetObjekteSha() schnittstellen.Sha {
	return a.ObjekteSha
}

func (a Sku) GetAkteSha() schnittstellen.Sha {
	return a.AkteSha
}

func (a Sku) Less(b Sku) (ok bool) {
	if a.Tai.Less(b.Tai) {
		ok = true
		return
	}

	return
}

func (a Sku) EqualsAny(b any) (ok bool) {
	return values.Equals(a, b)
}

func (a Sku) Equals(b Sku) (ok bool) {
	if a != b {
		return false
	}

	return true
}

func (s Sku) String() string {
	return fmt.Sprintf(
		"%s %s %s %s %s",
		s.Tai,
		s.Gattung,
		s.Kennung,
		s.ObjekteSha,
		s.AkteSha,
	)
}
