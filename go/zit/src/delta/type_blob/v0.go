package type_blob

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/etikett_rule"
	"code.linenisgreat.com/zit/go/zit/src/alfa/reset"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
)

type V0 struct {
	InlineAkte    bool                        `toml:"inline-akte,omitempty"`
	Archived      bool                        `toml:"archived,omitempty"`
	FileExtension string                      `toml:"file-extension,omitempty"`
	ExecCommand   *script_config.ScriptConfig `toml:"exec-command,omitempty"`
	VimSyntaxType string                      `toml:"vim-syntax-type"`
	// TODO-P4 rename to uti-groups
	FormatterUTIGroups map[string]FormatterUTIGroup              `toml:"formatter-uti-groups"`
	Formatters         map[string]script_config.WithOutputFormat `toml:"formatters,omitempty"`
	Actions            map[string]script_config.ScriptConfig     `toml:"actions,omitempty"`
	EtikettenRules     map[string]etikett_rule.Rule              `toml:"etiketten-rules,omitempty"` // deprecated
	Hooks              interface{}                               `toml:"hooks"`
}

func (a *V0) Reset() {
	a.Archived = false
	a.InlineAkte = false
	a.FileExtension = ""
	a.ExecCommand = nil
	a.VimSyntaxType = ""

	a.FormatterUTIGroups = reset.Map(a.FormatterUTIGroups)
	a.Formatters = reset.Map(a.Formatters)
	a.Actions = reset.Map(a.Actions)
	a.EtikettenRules = reset.Map(a.EtikettenRules)
	a.Hooks = nil
}
