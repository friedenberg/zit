package konfig

// TODO-P4 rename to Akte
type Toml struct {
	FileExtensions FileExtensions          `toml:"file-extensions"`
	RemoteScripts  map[string]RemoteScript `toml:"remote-scripts"`
	Recipients     []string                `toml:"recipients"`
}

type FileExtensions struct {
	Zettel   string `toml:"zettel"`
	Organize string `toml:"organize"`
	Typ      string `toml:"typ"`
	Etikett  string `toml:"etikett"`
}
