package typ

import "github.com/friedenberg/zit/src/bravo/script_config"

type Action struct {
	Description string `toml:"description"`
	//Must be embedded for toml formatting purposes
	script_config.ScriptConfig
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
