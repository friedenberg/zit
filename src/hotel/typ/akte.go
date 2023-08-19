package typ

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/etikett_rule"
	"github.com/friedenberg/zit/src/alfa/reset"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/script_config"
)

type Akte struct {
	InlineAkte    bool                        `toml:"inline-akte,omitempty"`
	Archived      bool                        `toml:"archived,omitempty"`
	FileExtension string                      `toml:"file-extension,omitempty"`
	ExecCommand   *script_config.ScriptConfig `toml:"exec-command,omitempty"`
	VimSyntaxType string                      `toml:"vim-syntax-type"`
	// TODO-P4 rename to uti-groups
	FormatterUTIGroups map[string]FormatterUTIGroup              `toml:"formatter-uti-groups"`
	Formatters         map[string]script_config.WithOutputFormat `toml:"formatters,omitempty"`
	Actions            map[string]script_config.ScriptConfig     `toml:"actions,omitempty"`
	EtikettenRules     map[string]etikett_rule.Rule              `toml:"etiketten-rules,omitempty"`
}

func (a Akte) GetGattung() schnittstellen.GattungLike {
	return gattung.Typ
}

func (a *Akte) Reset() {
	a.Archived = false
	a.InlineAkte = true
	a.FileExtension = ""
	a.ExecCommand = nil
	a.VimSyntaxType = ""

	a.FormatterUTIGroups = reset.Map(a.FormatterUTIGroups)
	a.Formatters = reset.Map(a.Formatters)
	a.Actions = reset.Map(a.Actions)
	a.EtikettenRules = reset.Map(a.EtikettenRules)
}

func (a *Akte) ResetWith(b Akte) {
	a.InlineAkte = b.InlineAkte
	a.Archived = b.Archived
	a.FileExtension = b.FileExtension
	a.ExecCommand = b.ExecCommand
	a.VimSyntaxType = b.VimSyntaxType

	errors.TodoP1("copy instead of assign")
	a.FormatterUTIGroups = b.FormatterUTIGroups
	a.Formatters = b.Formatters
	a.Actions = b.Actions
	a.EtikettenRules = b.EtikettenRules
}

func (a Akte) Equals(b Akte) bool {
	if a.Archived != b.Archived {
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

	if len(a.FormatterUTIGroups) != len(b.FormatterUTIGroups) {
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

		if !v.Equals(v1) {
			return false
		}
	}

	for k, v := range a.FormatterUTIGroups {
		v1, ok := b.FormatterUTIGroups[k]

		if !ok {
			return false
		}

		if !v.Equals(v1) {
			return false
		}
	}

	for k, v := range a.Formatters {
		v1, ok := b.Formatters[k]

		if !ok {
			return false
		}

		if !v.Equals(v1) {
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

	if a.VimSyntaxType != b.VimSyntaxType {
		return false
	}

	return true
}

func (a *Akte) Apply(b Akte) {
	a.Archived = b.Archived
	a.InlineAkte = b.InlineAkte
	a.FileExtension = b.FileExtension

	if len(b.Actions) > 0 {
		a.Actions = b.Actions
	}

	// if b.Description != "" {
	// 	a.Description = collections.MakeStringValue(b.Description)
	// }

	if len(b.Actions) > 0 {
		a.Actions = b.Actions
	}

	if len(b.Formatters) > 0 {
		a.Formatters = b.Formatters
	}

	if len(b.FormatterUTIGroups) > 0 {
		a.FormatterUTIGroups = b.FormatterUTIGroups
	}

	if b.ExecCommand != nil {
		a.ExecCommand = b.ExecCommand
	}

	if len(b.EtikettenRules) > 0 {
		a.EtikettenRules = b.EtikettenRules
	}

	if len(b.VimSyntaxType) > 0 {
		a.VimSyntaxType = b.VimSyntaxType
	}
}

func (a *Akte) Merge(b Akte) {
	if b.InlineAkte {
		a.InlineAkte = true
	}

	if b.FileExtension != "" {
		a.FileExtension = b.FileExtension
	}

	if b.ExecCommand != nil {
		a.ExecCommand.Merge(*b.ExecCommand)
	}

	for k, v := range b.EtikettenRules {
		a.EtikettenRules[k] = v
	}

	for k, v := range b.Actions {
		sc, ok := a.Actions[k]

		if !ok {
			sc = v
		} else {
			sc.Merge(v)
		}

		a.Actions[k] = sc
	}

	for k, v := range b.Formatters {
		sc, ok := a.Formatters[k]

		if !ok {
			sc = v
		} else {
			sc.Merge(v)
		}

		a.Formatters[k] = sc
	}

	for k, v := range b.FormatterUTIGroups {
		sc, ok := a.FormatterUTIGroups[k]

		if !ok {
			sc = v
		} else {
			sc.Merge(v)
		}

		a.FormatterUTIGroups[k] = sc
	}

	if len(b.VimSyntaxType) > 0 {
		a.VimSyntaxType = b.VimSyntaxType
	}
}
