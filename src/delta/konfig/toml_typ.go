package konfig

type KonfigTyp struct {
	InlineAkte     bool                        `toml:"inline-akte" default:"true"`
	FileExtension  string                      `toml:"file-extension"`
	ExecCommand    *ScriptConfig               `toml:"exec-command"`
	Actions        map[string]*KonfigTypAction `toml:"actions"`
	EtikettenRules map[string]EtikettRule      `toml:"etiketten-rules"`
}

func (a *KonfigTyp) Equals(b *KonfigTyp) bool {
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

	for k, v := range a.Actions {
		v1, ok := b.Actions[k]

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

	return true
}

type KonfigTypAction struct {
	ScriptConfig
	Description string `toml:"description"`
}

func (a *KonfigTypAction) Equals(b *KonfigTypAction) bool {
	if !a.ScriptConfig.Equals(&b.ScriptConfig) {
		return false
	}

	if a.Description != b.Description {
		return false
	}

	return true
}

func (kta *KonfigTypAction) Merge(kta2 *KonfigTypAction) {
	if kta2.Description != "" {
		kta.Description = kta2.Description
	}

	kta.ScriptConfig.Merge(&kta.ScriptConfig)
}

type EtikettRule struct {
	GoldenChild EtikettRuleGoldenChild `toml:"golden-child"`
}
