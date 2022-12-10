package etikett_rule

import "strings"

type RuleGoldenChild int

const (
	RuleGoldenChildUnset = RuleGoldenChild(iota)
	RuleGoldenChildLowest
	RuleGoldenChildHighest
)

func (t *RuleGoldenChild) Set(v string) (err error) {
	v = strings.ToLower(strings.TrimSpace(v))

	switch v {
	default:
		*t = RuleGoldenChildUnset

	case "lowest":
		*t = RuleGoldenChildLowest

	case "highest":
		*t = RuleGoldenChildLowest
	}

	return
}

func (t RuleGoldenChild) String() string {
	switch t {
	default:
		return ""

	case RuleGoldenChildLowest:
		return "lowest"

	case RuleGoldenChildHighest:
		return "highest"
	}
}

func (t RuleGoldenChild) MarshalText() (text []byte, err error) {
	text = []byte(t.String())

	return
}

func (t *RuleGoldenChild) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		return
	}

	return
}
