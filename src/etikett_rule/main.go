package etikett_rule

type Rule struct {
	GoldenChild RuleGoldenChild `toml:"golden-child,omitempty"`
}
