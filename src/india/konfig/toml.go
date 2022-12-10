package konfig

import "github.com/friedenberg/zit/src/etikett"

// TODO-P4 rename to Akte
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

type KonfigTag = etikett.Akte
