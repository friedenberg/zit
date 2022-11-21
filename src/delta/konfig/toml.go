package konfig

type tomlKonfig struct {
	RemoteScripts map[string]RemoteScript `toml:"remote-scripts"`
	Tags          map[string]KonfigTag    `toml:"tags"`
	Typen         map[string]KonfigTyp    `toml:"typen"`
	Recipients    []string                `toml:"recipients"`
}

type KonfigTag struct {
	AddToNewZettels bool `toml:"add-to-new-zettels"`
	Hide            bool `toml:"hide"`
}

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
