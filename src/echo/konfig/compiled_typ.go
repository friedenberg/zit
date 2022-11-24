package konfig

import (
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
)

var (
	typExpander kennung.Expander
)

func init() {
	typExpander = kennung.MakeExpanderRight(`-`)
}

type compiledTyp struct {
	Name           collections.StringValue
	InlineAkte     bool
	FileExtension  string
	ExecCommand    *ScriptConfig
	Actions        map[string]*KonfigTypAction
	EtikettenRules map[string]EtikettRule
}

func makeCompiledTyp(n string) *compiledTyp {
	return &compiledTyp{
		Name:           collections.MakeStringValue(n),
		ExecCommand:    &ScriptConfig{},
		Actions:        make(map[string]*KonfigTypAction),
		EtikettenRules: make(map[string]EtikettRule),
	}
}

func (ct *compiledTyp) Apply(kt KonfigTyp) {
	ct.InlineAkte = kt.InlineAkte
	ct.FileExtension = kt.FileExtension

	// if kt.Description != "" {
	// 	ct.Description = collections.MakeStringValue(kt.Description)
	// }

	if len(kt.Actions) > 0 {
		ct.Actions = kt.Actions
	}

	if kt.ExecCommand != nil {
		ct.ExecCommand = kt.ExecCommand
	}

	if len(kt.EtikettenRules) > 0 {
		ct.EtikettenRules = kt.EtikettenRules
	}

}

func (ct *compiledTyp) Merge(ct2 *compiledTyp) {
	if ct2.InlineAkte {
		ct.InlineAkte = true
	}

	if ct2.FileExtension != "" {
		ct.FileExtension = ct2.FileExtension
	}

	ct.ExecCommand.Merge(ct2.ExecCommand)

	for k, v := range ct2.EtikettenRules {
		ct.EtikettenRules[k] = v
	}

	for k, v := range ct2.Actions {
		sc, ok := ct.Actions[k]

		if !ok {
			sc = v
		} else {
			sc.Merge(v)
		}

		ct.Actions[k] = sc
	}
}

func (ct *compiledTyp) ApplyExpanded(c Compiled) {
	expandedActual := c.GetSortedTypenExpanded(ct.Name.String())

	for _, ex := range expandedActual {
		ct.Merge(ex)
	}
}
