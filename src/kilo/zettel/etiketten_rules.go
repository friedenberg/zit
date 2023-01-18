package zettel

import (
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/etikett_rule"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/juliett/konfig"
)

func (z *Objekte) ApplyKonfig(k konfig.Compiled) (err error) {
	errors.TodoP4("move this somewhere more appropriate")

	normalized := kennung.WithRemovedCommonPrefixes(z.Etiketten)
	z.Etiketten = normalized

	tk := k.GetTyp(z.Typ)

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
	mode etikett_rule.RuleGoldenChild,
) (err error) {
	if z.Etiketten.Len() == 0 {
		return
	}

	switch mode {
	case etikett_rule.RuleGoldenChildUnset:
		return
	}

	mes := z.Etiketten.MutableCopy()

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
	z.Etiketten = mes.Copy()

	return
}
