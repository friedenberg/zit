package metadatei

import (
	"flag"
	"io"
	"sort"
	"strings"

	"github.com/friedenberg/zit/src/alfa/etikett_rule"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/collections2"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
)

var BoundaryStringValue values.String

const (
	Boundary = "---"
)

func init() {
	BoundaryStringValue = values.MakeString(Boundary)
}

type MetadateiWriterTo interface {
	io.WriterTo
	HasMetadateiContent() bool
}

type Metadatei struct {
	// StoreVersion values.Int
	// Kasten
	// Domain
	AkteSha     sha.Sha
	Bezeichnung bezeichnung.Bezeichnung
	Etiketten   kennung.EtikettSet
	// Gattung     gattung.Gattung
	Typ kennung.Typ
	Tai kennung.Tai
}

func (m Metadatei) GetMetadatei() Metadatei {
	return m
}

func (a *Metadatei) SetMetadatei(b Metadatei) {
	*a = b
}

func (m *Metadatei) AddToFlagSet(f *flag.FlagSet) {
	f.Var(
		&m.Bezeichnung,
		"bezeichnung",
		"the Bezeichnung to use for created or updated Zettelen",
	)

	mes := m.GetEtiketten().CloneMutableSetPtrLike()

	fes := collections2.MakeFlagCommasFromExisting[kennung.Etikett](
		collections2.SetterPolicyAppend,
		mes,
	)

	f.Var(
		fes,
		"etiketten",
		"the Etiketten to use for created or updated Zettelen",
	)

	m.Etiketten = mes

	// TODO-P1 add typ
}

func (z Metadatei) UserInputIsEmpty() bool {
	if !z.Bezeichnung.IsEmpty() {
		return false
	}

	if z.Etiketten != nil && z.Etiketten.Len() > 0 {
		return false
	}

	if !kennung.IsEmpty(z.Typ) {
		return false
	}

	return true
}

func (z Metadatei) IsEmpty() bool {
	if !z.AkteSha.IsNull() {
		return false
	}

	if !z.UserInputIsEmpty() {
		return false
	}

	if !z.Tai.IsZero() {
		return false
	}

	return true
}

func (z *Metadatei) SetBezeichnung(b bezeichnung.Bezeichnung) {
	z.Bezeichnung = b
}

func (z *Metadatei) SetEtiketten(e kennung.EtikettSet) {
	z.Etiketten = e
}

func (z *Metadatei) SetTyp(t kennung.Typ) {
	z.Typ = t
}

func (z Metadatei) GetBezeichnung() bezeichnung.Bezeichnung {
	return z.Bezeichnung
}

func (z Metadatei) GetEtiketten() kennung.EtikettSet {
	if z.Etiketten == nil {
		return kennung.MakeEtikettSet()
	}

	return z.Etiketten
}

func (z Metadatei) GetTyp() kennung.Typ {
	return z.Typ
}

func (z Metadatei) GetTai() kennung.Tai {
	return z.Tai
}

func (pz Metadatei) EqualsSansTai(z1 Metadatei) bool {
	if !pz.AkteSha.Equals(z1.AkteSha) {
		return false
	}

	if !pz.Typ.Equals(z1.Typ) {
		return false
	}

	switch {
	case pz.Etiketten == nil && z1.Etiketten == nil:
	// pass
	case pz.Etiketten == nil:
		return z1.Etiketten.Len() == 0
	case z1.Etiketten == nil:
		return pz.Etiketten.Len() == 0
	case !pz.Etiketten.EqualsSetLike(z1.Etiketten):
		return false
	}

	if !pz.Bezeichnung.Equals(z1.Bezeichnung) {
		return false
	}

	return true
}

func (pz Metadatei) Equals(z1 Metadatei) bool {
	if !pz.EqualsSansTai(z1) {
		return false
	}

	if !pz.Tai.Equals(z1.Tai) {
		return false
	}

	return true
}

func (z *Metadatei) Reset() {
	z.AkteSha.Reset()
	z.Bezeichnung.Reset()
	z.Etiketten = kennung.MakeEtikettSet()
	z.Typ = kennung.Typ{}
	// z.Gattung = gattung.Unknown
	z.Tai.Reset()
}

func (z *Metadatei) ResetWith(z1 Metadatei) {
	z.AkteSha = z1.AkteSha
	z.Bezeichnung = z1.Bezeichnung

	if z1.Etiketten == nil {
		z.Etiketten = kennung.MakeEtikettSet()
	} else {
		z.Etiketten = z1.Etiketten.CloneSetPtrLike()
	}

	z.Typ = z1.Typ
	// z.Gattung = z1.Gattung
	z.Tai = z1.Tai
}

func (z Metadatei) Description() (d string) {
	d = z.Bezeichnung.String()

	if strings.TrimSpace(d) == "" {
		d = iter.StringCommaSeparated[kennung.Etikett](z.Etiketten)
	}

	return
}

func (z *Metadatei) ApplyGoldenChild(
	e kennung.Etikett,
	mode etikett_rule.RuleGoldenChild,
) (err error) {
	if z.Etiketten.Len() == 0 {
		return
	}

	switch mode {
	case etikett_rule.RuleGoldenChildUnset:
		return
	}

	mes := z.Etiketten.CloneMutableSetPtrLike()

	prefixes := kennung.Withdraw(mes, e).Elements()

	if len(prefixes) == 0 {
		return
	}

	var sortFunc func(i, j int) bool

	switch mode {
	case etikett_rule.RuleGoldenChildLowest:
		sortFunc = func(i, j int) bool { return kennung.Less(prefixes[j], prefixes[i]) }

	case etikett_rule.RuleGoldenChildHighest:
		sortFunc = func(i, j int) bool { return kennung.Less(prefixes[i], prefixes[j]) }
	}

	sort.Slice(prefixes, sortFunc)

	mes.Add(prefixes[0])
	z.Etiketten = mes.CloneSetPtrLike()

	return
}
