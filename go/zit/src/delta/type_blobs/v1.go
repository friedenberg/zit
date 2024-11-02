package type_blobs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/reset"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
)

type TomlV1 struct {
	Binary        bool                                      `toml:"binary,omitempty"`
	FileExtension string                                    `toml:"file-extension,omitempty"`
	ExecCommand   *script_config.ScriptConfig               `toml:"exec-command,omitempty"`
	VimSyntaxType string                                    `toml:"vim-syntax-type"`
	UTIGroups     map[string]FormatterUTIGroup              `toml:"uti-groups"`
	Formatters    map[string]script_config.WithOutputFormat `toml:"formatters,omitempty"`
	Hooks         interface{}                               `toml:"hooks"`
}

func (a *TomlV1) Reset() {
	a.Binary = false
	a.FileExtension = ""
	a.ExecCommand = nil
	a.VimSyntaxType = ""

	a.UTIGroups = reset.Map(a.UTIGroups)
	a.Formatters = reset.Map(a.Formatters)
	a.Hooks = nil
}

func (a *TomlV1) GetBinary() bool {
  return a.Binary
}

func (a *TomlV1) GetFileExtension() string {
  return a.FileExtension
}
