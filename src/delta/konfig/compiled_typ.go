package konfig

import (
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/charlie/string_expansion"
)

var (
	typExpander string_expansion.Expander[collections.StringValue, *collections.StringValue]
)

func init() {
	typExpander = string_expansion.MakeExpanderRight[collections.StringValue, *collections.StringValue](`-`)
}

type compiledTyp struct {
	Name           collections.StringValue
	InlineAkte     bool
	FileExtension  string
	FormatScript   *ScriptConfig
	ExecCommand    *ScriptConfig
	Actions        map[string]*ScriptConfig
	EtikettenRules map[string]EtikettRule
}

func makeCompiledTyp(n string) *compiledTyp {
	return &compiledTyp{
		Name:           collections.MakeStringValue(n),
		FormatScript:   &ScriptConfig{},
		ExecCommand:    &ScriptConfig{},
		Actions:        make(map[string]*ScriptConfig),
		EtikettenRules: make(map[string]EtikettRule),
	}
}

func (ct *compiledTyp) Apply(kt KonfigTyp) {
	ct.InlineAkte = kt.InlineAkte
	ct.FileExtension = kt.FileExtension

	if len(kt.Actions) > 0 {
		ct.Actions = kt.Actions
	}

	if kt.FormatScript != nil {
		ct.FormatScript = kt.FormatScript
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

	ct.FormatScript.Merge(ct2.FormatScript)
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
