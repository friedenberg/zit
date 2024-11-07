package options_print

type BoxV0 struct {
	PrintIncludeDescription bool `toml:"print-include-description"`
	PrintTime               bool `toml:"print-time"`
	PrintTai                bool `toml:"-"`
	PrintTagsAlways         bool `toml:"print-etiketten-always"`
	PrintEmptyShas          bool `toml:"print-empty-shas"`
	PrintIncludeTypes       bool `toml:"print-include-typen"`
	DescriptionInBox        bool `toml:"-"`
	ExcludeFields           bool `toml:"-"`
	PrintState              bool `toml:"-"`
}
