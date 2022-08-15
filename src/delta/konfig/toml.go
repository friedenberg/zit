package konfig

type Toml struct {
	RemoteScripts  map[string]RemoteScriptConfig `toml:"remote-scripts"`
	Tags           map[string]KonfigTag          `toml:"tags"`
	Clean          string                        `toml:"clean"`
	Smudge         string                        `toml:"smudge"`
}
