package konfig

type Toml struct {
	RemoteScripts map[string]RemoteScript `toml:"remote-scripts"`
	Tags          map[string]KonfigTag    `toml:"tags"`
	Clean         string                  `toml:"clean"`
	Smudge        string                  `toml:"smudge"`
	Typen         map[string]KonfigTyp    `toml:"typen"`
}
