package konfig

//TODO-P4 rename to Akte
type Toml struct {
	FileExtensions FileExtensions          `toml:"file-extensions"`
	RemoteScripts  map[string]RemoteScript `toml:"remote-scripts"`
	Tags           map[string]KonfigTag    `toml:"tags"`
	Recipients     []string                `toml:"recipients"`
}

type FileExtensions struct {
	Zettel   string `toml:"zettel"`
	Organize string `toml:"organize"`
	Typ      string `toml:"typ"`
}

//TODO-P2 move to etikett package
type KonfigTag struct {
	AddToNewZettels bool `toml:"add-to-new-zettels"`
	Hide            bool `toml:"hide"`
}
