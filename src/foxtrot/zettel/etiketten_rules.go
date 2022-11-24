package zettel

import (
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/kennung"
	"github.com/friedenberg/zit/src/delta/konfig"
)

// TODO-P2 move this to somewhere more appropriate
func (z *Zettel) ApplyKonfig(k konfig.Konfig) (err error) {
	normalized := kennung.WithRemovedCommonPrefixes(z.Etiketten)
	z.Etiketten = normalized

	tk := k.Compiled.GetTyp(z.Typ.String())

	if tk == nil {
		return
	}

	for e, r := range tk.EtikettenRules {
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

func (z *Zettel) applyGoldenChild(
	e kennung.Etikett,
	mode konfig.EtikettRuleGoldenChild,
) (err error) {
	if z.Etiketten.Len() == 0 {
		return
	}

	switch mode {
	case konfig.EtikettRuleGoldenChildUnset:
		return
	}

	mes := z.Etiketten.MutableCopy()

	prefixes := kennung.Withdraw(mes, e).Elements()

	if len(prefixes) == 0 {
		return
	}

	var sortFunc func(i, j int) bool

	switch mode {
	case konfig.EtikettRuleGoldenChildLowest:
		sortFunc = func(i, j int) bool { return prefixes[j].Less(prefixes[i]) }

	case konfig.EtikettRuleGoldenChildHighest:
		sortFunc = func(i, j int) bool { return prefixes[i].Less(prefixes[j]) }
	}

	sort.Slice(prefixes, sortFunc)

	mes.Add(prefixes[0])
	z.Etiketten = mes.Copy()

	return
}
