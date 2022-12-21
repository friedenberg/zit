package typ

import (
	"github.com/friedenberg/zit/src/alfa/etikett_rule"
	"github.com/friedenberg/zit/src/bravo/script_config"
)

type Akte struct {
	InlineAkte     bool                                         `toml:"inline-akte,omitempty"`
	FileExtension  string                                       `toml:"file-extension,omitempty"`
	ExecCommand    *script_config.ScriptConfig                  `toml:"exec-command,omitempty"`
	Formatters     map[string]script_config.ScriptConfigWithUTI `toml:"formatters,omitempty"`
	Actions        map[string]script_config.ScriptConfig        `toml:"actions,omitempty"`
	EtikettenRules map[string]etikett_rule.Rule                 `toml:"etiketten-rules,omitempty"`
}

func (a *Akte) Reset(b *Akte) {
	panic("not implemented")
	// if b == nil {
	// } else {
	// }
}

func (a *Akte) Equals(b *Akte) bool {
	if b == nil || a == nil {
		return false
	}

	if a.InlineAkte != b.InlineAkte {
		return false
	}

	if a.FileExtension != b.FileExtension {
		return false
	}

	if a.ExecCommand != b.ExecCommand {
		return false
	}

	if len(a.Actions) != len(b.Actions) {
		return false
	}

	if len(a.Formatters) != len(b.Formatters) {
		return false
	}

	for k, v := range a.Actions {
		v1, ok := b.Actions[k]

		if !ok {
			return false
		}

		if !v.Equals(&v1) {
			return false
		}
	}

	for k, v := range a.Formatters {
		v1, ok := b.Formatters[k]

		if !ok {
			return false
		}

		if !v.Equals(&v1) {
			return false
		}
	}

	if len(a.EtikettenRules) != len(b.EtikettenRules) {
		return false
	}

	for k, v := range a.EtikettenRules {
		v1, ok := b.EtikettenRules[k]

		if !ok {
			return false
		}

		if v != v1 {
			return false
		}
	}

	return true
}

func (ct *Akte) Apply(kt *Akte) {
	ct.InlineAkte = kt.InlineAkte
	ct.FileExtension = kt.FileExtension

	// if kt.Description != "" {
	// 	ct.Description = collections.MakeStringValue(kt.Description)
	// }

	if len(kt.Actions) > 0 {
		ct.Actions = kt.Actions
	}

	if len(kt.Formatters) > 0 {
		ct.Formatters = kt.Formatters
	}

	if kt.ExecCommand != nil {
		ct.ExecCommand = kt.ExecCommand
	}

	if len(kt.EtikettenRules) > 0 {
		ct.EtikettenRules = kt.EtikettenRules
	}
}

func (ct *Akte) Merge(ct2 *Akte) {
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
			sc.Merge(&v)
		}

		ct.Actions[k] = sc
	}

	for k, v := range ct2.Formatters {
		sc, ok := ct.Formatters[k]

		if !ok {
			sc = v
		} else {
			sc.Merge(&v)
		}

		ct.Formatters[k] = sc
	}
}
