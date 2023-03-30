package erworben

import (
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type Akte struct {
	DefaultTyp     kennung.Typ                           `toml:"default-typ"`
	FileExtensions FileExtensions                        `toml:"file-extensions"`
	RemoteScripts  map[string]RemoteScript               `toml:"remote-scripts"`
	Recipients     []string                              `toml:"recipients"`
	Actions        map[string]script_config.ScriptConfig `toml:"actions,omitempty"`
}

func (a *Akte) Reset() {
	a.FileExtensions.Reset()
	a.DefaultTyp = kennung.Typ{}
	a.RemoteScripts = make(map[string]RemoteScript)
	// TODO-P4 should reuse
	a.Recipients = make([]string, 0)
	a.Actions = make(map[string]script_config.ScriptConfig)
}

func (a *Akte) ResetWith(b Akte) {
	a.FileExtensions.Reset()
	// TODO-P4 should copy
	a.DefaultTyp = b.DefaultTyp
	a.RemoteScripts = b.RemoteScripts
	a.Recipients = b.Recipients
	a.Actions = b.Actions
}
