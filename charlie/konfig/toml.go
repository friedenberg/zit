package konfig

type Toml struct {
	//TODO
	DoNotLockFiles bool                          `toml:"do-not-lock-files"`
	RemoteScripts  map[string]RemoteScriptConfig `toml:"remote-scripts"`
	Tags           map[string]KonfigTag          `toml:"tags"`
}
