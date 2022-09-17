package zettel

import (
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/konfig"
)

func (z *Zettel) ApplyTypKonfig(tk konfig.KonfigTyp) (err error) {
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

	prefixes := z.Etiketten.Withdraw(e).Etiketten()

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

	z.Etiketten.Add(prefixes[0])

	return
}
