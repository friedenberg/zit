package typ_toml

type EtikettRule struct {
	GoldenChild EtikettRuleGoldenChild `toml:"golden-child,omitempty"`
}
