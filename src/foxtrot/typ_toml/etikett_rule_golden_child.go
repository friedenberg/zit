package typ_toml

import "strings"

type EtikettRuleGoldenChild int

const (
	EtikettRuleGoldenChildUnset = EtikettRuleGoldenChild(iota)
	EtikettRuleGoldenChildLowest
	EtikettRuleGoldenChildHighest
)

func (t *EtikettRuleGoldenChild) Set(v string) (err error) {
	v = strings.ToLower(strings.TrimSpace(v))

	switch v {
	default:
		*t = EtikettRuleGoldenChildUnset

	case "lowest":
		*t = EtikettRuleGoldenChildLowest

	case "highest":
		*t = EtikettRuleGoldenChildLowest
	}

	return
}

func (t EtikettRuleGoldenChild) String() string {
	switch t {
	default:
		return ""

	case EtikettRuleGoldenChildLowest:
		return "lowest"

	case EtikettRuleGoldenChildHighest:
		return "highest"
	}
}

func (t EtikettRuleGoldenChild) MarshalText() (text []byte, err error) {
	text = []byte(t.String())

	return
}

func (t *EtikettRuleGoldenChild) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		return
	}

	return
}
