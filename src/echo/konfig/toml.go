package konfig

import "github.com/friedenberg/zit/src/delta/typ_toml"

type tomlKonfig struct {
	RemoteScripts map[string]RemoteScript `toml:"remote-scripts"`
	Tags          map[string]KonfigTag    `toml:"tags"`
	Typen         map[string]typ_toml.Typ `toml:"typen"`
	Recipients    []string                `toml:"recipients"`
}

type KonfigTag struct {
	AddToNewZettels bool `toml:"add-to-new-zettels"`
	Hide            bool `toml:"hide"`
}
