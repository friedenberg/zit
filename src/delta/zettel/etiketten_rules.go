package zettel

import (
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/konfig"
)

//TODO move this to somewhere more appropriate
func (z *Zettel) ApplyKonfig(k konfig.Konfig) (err error) {
	normalized := etikett.WithRemovedCommonPrefixes(z.Etiketten)
	z.Etiketten = normalized

	var tk konfig.KonfigTyp
	ok := false

	if tk, ok = k.Typen[z.Typ.String()]; !ok {
		return
	}

	for e, r := range tk.EtikettenRules {
		e1 := etikett.Make(e)

		if err = z.applyGoldenChild(e1, r.GoldenChild); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (z *Zettel) applyGoldenChild(
	e etikett.Etikett,
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

	prefixes := etikett.Withdraw(mes, e).Elements()

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
