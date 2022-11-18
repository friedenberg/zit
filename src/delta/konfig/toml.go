package konfig

type toml struct {
	RemoteScripts map[string]RemoteScript `toml:"remote-scripts"`
	Tags          map[string]KonfigTag    `toml:"tags"`
	Clean         string                  `toml:"clean"`
	Smudge        string                  `toml:"smudge"`
	Typen         map[string]KonfigTyp    `toml:"typen"`
	Recipients    []string                `toml:"recipients"`
}

type KonfigTag struct {
	AddToNewZettels bool `toml:"add-to-new-zettels"`
	Hide            bool `toml:"hide"`
}

type KonfigTyp struct {
	FormatScript   ScriptConfig            `toml:"format-script"`
	InlineAkte     bool                    `toml:"inline-akte" default:"true"`
	Actions        map[string]ScriptConfig `toml:"actions"`
	ExecCommand    ScriptConfig            `toml:"exec-command"`
	EtikettenRules map[string]EtikettRule  `toml:"etiketten-rules"`
	FileExtension  string                  `toml:"file-extension"`
}

type EtikettRule struct {
	GoldenChild EtikettRuleGoldenChild `toml:"golden-child"`
}
