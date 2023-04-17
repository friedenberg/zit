package metadatei

import (
	"flag"
	"io"
	"sort"
	"strings"

	"github.com/friedenberg/zit/src/alfa/etikett_rule"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
)

const (
	Boundary = "---"
)

type MetadateiWriterTo interface {
	io.WriterTo
	HasMetadateiContent() bool
}

type Metadatei struct {
	AkteSha     sha.Sha
	Bezeichnung bezeichnung.Bezeichnung
	Etiketten   kennung.EtikettSet
	Typ         kennung.Typ
}

func (m *Metadatei) AddToFlagSet(f *flag.FlagSet) {
	f.Var(&m.Bezeichnung, "bezeichnung", "the Bezeichnung to use for created or updated Zettelen")
	f.Var(
		collections.MakeFlagCommasFromExisting(
			collections.SetterPolicyAppend,
			&m.Etiketten,
		),
		"etiketten",
		"the Etiketten to use for created or updated Zttelen",
	)
	// TODO add typ
}

func (z Metadatei) IsEmpty() bool {
	if !z.AkteSha.IsNull() {
		return false
	}

	if !z.Bezeichnung.IsEmpty() {
		return false
	}

	if z.Etiketten.Len() > 0 {
		return false
	}

	if !z.Typ.IsEmpty() {
		return false
	}

	return true
}

func (z *Metadatei) SetBezeichnung(b bezeichnung.Bezeichnung) {
	z.Bezeichnung = b
}

func (z *Metadatei) SetEtiketten(e schnittstellen.Set[kennung.Etikett]) {
	z.Etiketten = e
}

func (z *Metadatei) SetTyp(t kennung.Typ) {
	z.Typ = t
}

func (z Metadatei) GetBezeichnung() bezeichnung.Bezeichnung {
	return z.Bezeichnung
}

func (z Metadatei) GetEtiketten() schnittstellen.Set[kennung.Etikett] {
	return z.Etiketten.ImmutableClone()
}

func (z Metadatei) GetTyp() kennung.Typ {
	return z.Typ
}

func (pz Metadatei) Equals(z1 Metadatei) bool {
	if !pz.AkteSha.Equals(z1.AkteSha) {
		return false
	}

	if !pz.Typ.Equals(z1.Typ) {
		return false
	}

	if !pz.Etiketten.Equals(z1.Etiketten) {
		return false
	}

	if !pz.Bezeichnung.Equals(z1.Bezeichnung) {
		return false
	}

	return true
}

func (z *Metadatei) Reset() {
	z.AkteSha.Reset()
	z.Bezeichnung.Reset()
	z.Etiketten = kennung.MakeEtikettSet()
	z.Typ = kennung.Typ{}
}

func (z *Metadatei) ResetWith(z1 Metadatei) {
	z.AkteSha = z1.AkteSha
	z.Bezeichnung = z1.Bezeichnung
	z.Etiketten = z1.Etiketten.ImmutableClone()
	z.Typ = z1.Typ
}

func (z Metadatei) Description() (d string) {
	d = z.Bezeichnung.String()

	if strings.TrimSpace(d) == "" {
		d = collections.StringCommaSeparated[kennung.Etikett](z.Etiketten)
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

	mes := z.Etiketten.MutableClone()

	prefixes := kennung.Withdraw(mes, e).Elements()

	if len(prefixes) == 0 {
		return
	}

	var sortFunc func(i, j int) bool

	switch mode {
	case etikett_rule.RuleGoldenChildLowest:
		sortFunc = func(i, j int) bool { return prefixes[j].Less(prefixes[i]) }

	case etikett_rule.RuleGoldenChildHighest:
		sortFunc = func(i, j int) bool { return prefixes[i].Less(prefixes[j]) }
	}

	sort.Slice(prefixes, sortFunc)

	mes.Add(prefixes[0])
	z.Etiketten = mes.ImmutableClone()

	return
}
