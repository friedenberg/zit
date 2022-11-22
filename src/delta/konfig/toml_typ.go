package konfig

type KonfigTypAction struct {
	ScriptConfig
	Description string `toml:"description"`
}

func (kta *KonfigTypAction) Merge(kta2 *KonfigTypAction) {
	if kta2.Description != "" {
		kta.Description = kta2.Description
	}

	kta.ScriptConfig.Merge(&kta.ScriptConfig)
}

type KonfigTyp struct {
	InlineAkte     bool                        `toml:"inline-akte" default:"true"`
	FileExtension  string                      `toml:"file-extension"`
	ExecCommand    *ScriptConfig               `toml:"exec-command"`
	Actions        map[string]*KonfigTypAction `toml:"actions"`
	EtikettenRules map[string]EtikettRule      `toml:"etiketten-rules"`
}

type EtikettRule struct {
	GoldenChild EtikettRuleGoldenChild `toml:"golden-child"`
}
