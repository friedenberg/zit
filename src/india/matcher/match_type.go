package matcher

type matchType int

const (
	matchTypeEmpty = matchType(iota)
	matchTypeStrictHinweis
	matchTypeCompound
)

func (mt *matchType) AddExact(m Matcher) {
	switch *mt {
	case matchTypeCompound:
		// noop
	default:
		*mt = matchTypeStrictHinweis
	}
}

func (mt *matchType) AddNonExact(m Matcher) {
	switch *mt {
	default:
		*mt = matchTypeCompound
	}
}
