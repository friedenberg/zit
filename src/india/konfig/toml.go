package konfig

import "github.com/friedenberg/zit/src/bravo/script_config"

// TODO-P4 rename to Akte
type Toml struct {
	FileExtensions FileExtensions                        `toml:"file-extensions"`
	RemoteScripts  map[string]RemoteScript               `toml:"remote-scripts"`
	Recipients     []string                              `toml:"recipients"`
	Actions        map[string]script_config.ScriptConfig `toml:"actions,omitempty"`
}

type FileExtensions struct {
	Zettel   string `toml:"zettel"`
	Organize string `toml:"organize"`
	Typ      string `toml:"typ"`
	Etikett  string `toml:"etikett"`
}
