package konfig

import "github.com/friedenberg/zit/src/bravo/script_config"

// TODO-P4 rename to Akte
type Toml struct {
	FileExtensions FileExtensions                        `toml:"file-extensions"`
	RemoteScripts  map[string]RemoteScript               `toml:"remote-scripts"`
	Recipients     []string                              `toml:"recipients"`
	Actions        map[string]script_config.ScriptConfig `toml:"actions,omitempty"`
	StoreVersion   string                                `toml:"store-version,omitempty"`
}

func (a *Toml) Reset(b *Toml) {
	if b == nil {
		a.FileExtensions.Reset(nil)
		a.RemoteScripts = make(map[string]RemoteScript)
		//TODO-P4 should reuse
		a.Recipients = make([]string, 0)
		a.Actions = make(map[string]script_config.ScriptConfig)
		a.StoreVersion = ""
	} else {
		a.FileExtensions.Reset(&b.FileExtensions)
		//TODO-P4 should copy
		a.RemoteScripts = b.RemoteScripts
		a.Recipients = b.Recipients
		a.Actions = b.Actions
		a.StoreVersion = b.StoreVersion
	}
}
