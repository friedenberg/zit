package konfig

type Toml struct {
	DoNotLockFiles bool                          `toml:"do-not-lock-files"`
	RemoteScripts  map[string]RemoteScriptConfig `toml:"remote-scripts"`
	Tags           map[string]KonfigTag          `toml:"tags"`
}
