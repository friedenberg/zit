package zettel

import (
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/typ_toml"
	"github.com/friedenberg/zit/src/hotel/konfig_compiled"
)

// TODO-P2 move this to somewhere more appropriate
func (z *Objekte) ApplyKonfig(k konfig_compiled.Compiled) (err error) {
	normalized := kennung.WithRemovedCommonPrefixes(z.Etiketten)
	z.Etiketten = normalized

	tk := k.GetTyp(z.Typ.String())

	if tk == nil {
		return
	}

	for e, r := range tk.Objekte.Akte.EtikettenRules {
		var e1 kennung.Etikett

		if e1, err = kennung.MakeEtikett(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = z.applyGoldenChild(e1, r.GoldenChild); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (z *Objekte) applyGoldenChild(
	e kennung.Etikett,
	mode typ_toml.EtikettRuleGoldenChild,
) (err error) {
	if z.Etiketten.Len() == 0 {
		return
	}

	switch mode {
	case typ_toml.EtikettRuleGoldenChildUnset:
		return
	}

	mes := z.Etiketten.MutableCopy()

	prefixes := kennung.Withdraw(mes, e).Elements()

	if len(prefixes) == 0 {
		return
	}

	var sortFunc func(i, j int) bool

	switch mode {
	case typ_toml.EtikettRuleGoldenChildLowest:
		sortFunc = func(i, j int) bool { return prefixes[j].Less(prefixes[i]) }

	case typ_toml.EtikettRuleGoldenChildHighest:
		sortFunc = func(i, j int) bool { return prefixes[i].Less(prefixes[j]) }
	}

	sort.Slice(prefixes, sortFunc)

	mes.Add(prefixes[0])
	z.Etiketten = mes.Copy()

	return
}
