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
	FormatScript   ScriptConfig
	InlineAkte     bool
	Actions        map[string]ScriptConfig
	ExecCommand    ScriptConfig
	EtikettenRules map[string]EtikettRule
	FileExtension  string
}

func (ct *compiledTyp) Apply(kt KonfigTyp) {
	ct.FormatScript = kt.FormatScript
	ct.InlineAkte = kt.InlineAkte
	ct.Actions = kt.Actions
	ct.ExecCommand = kt.ExecCommand
	ct.EtikettenRules = kt.EtikettenRules
	ct.FileExtension = kt.FileExtension
}

func (ct *compiledTyp) Merge(ct2 *compiledTyp) {
	ct.FormatScript.Merge(&ct2.FormatScript)

	if ct2.InlineAkte {
		ct.InlineAkte = true
	}

	if len(ct2.FormatScript.Shell) > 0 {
		ct.FormatScript.Shell = ct2.FormatScript.Shell
	}

	ct.ExecCommand.Merge(&ct2.ExecCommand)

	for k, v := range ct2.EtikettenRules {
		ct.EtikettenRules[k] = v
	}

	if ct2.FileExtension != "" {
		ct.FileExtension = ct2.FileExtension
	}

	for k, v := range ct2.Actions {
		sc, ok := ct.Actions[k]

		if !ok {
			sc = v
		} else {
			sc.Merge(&v)
		}

		ct.Actions[k] = v
	}
}

func (ct *compiledTyp) ApplyExpanded(c Compiled) {
	expandedActual := c.GetSortedTypenExpanded(ct.Name.String())

	for _, ex := range expandedActual {
		ct.Merge(ex)
	}
}
