package typ

import "github.com/friedenberg/zit/src/script_config"

type Action struct {
	ScriptConfig script_config.ScriptConfig
	Description  string `toml:"description"`
}

func (a *Action) Equals(b *Action) bool {
	if !a.ScriptConfig.Equals(&b.ScriptConfig) {
		return false
	}

	if a.Description != b.Description {
		return false
	}

	return true
}

func (kta *Action) Merge(kta2 *Action) {
	if kta2.Description != "" {
		kta.Description = kta2.Description
	}

	kta.ScriptConfig.Merge(&kta.ScriptConfig)
}
