package metadatei

import (
	"sort"
	"strings"

	"github.com/friedenberg/zit/src/alfa/etikett_rule"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
)

type Metadatei struct {
	Bezeichnung bezeichnung.Bezeichnung
	Etiketten   kennung.EtikettSet
}

func (z Metadatei) IsEmpty() bool {
	if !z.Bezeichnung.IsEmpty() {
		return false
	}

	if z.Etiketten.Len() > 0 {
		return false
	}

	return true
}

func (z Metadatei) GetEtiketten() schnittstellen.Set[kennung.Etikett] {
	return z.Etiketten.ImmutableClone()
}

func (z Metadatei) Equals(z1 Metadatei) bool {
	if !z.Bezeichnung.Equals(z1.Bezeichnung) {
		return false
	}

	if !z.Etiketten.Equals(z1.Etiketten) {
		return false
	}

	return true
}

func (z *Metadatei) Reset() {
	z.Bezeichnung.Reset()
	z.Etiketten = kennung.MakeEtikettSet()
}

func (z *Metadatei) ResetWith(z1 Metadatei) {
	z.Bezeichnung = z1.Bezeichnung
	z.Etiketten = z1.Etiketten.ImmutableClone()
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
