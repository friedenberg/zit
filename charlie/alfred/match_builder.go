package alfred

import (
	"strings"
)

type MatchBuilder struct {
	sb *strings.Builder
}

func NewMatchBuilder() MatchBuilder {
	return MatchBuilder{
		sb: &strings.Builder{},
	}
}

func (mb MatchBuilder) AddMatch(s string) {
	s1 := strings.Split(s, "_")

	for _, s2 := range s1 {
		mb.sb.WriteString(s2)
		mb.sb.WriteString(" ")
	}
}

func (mb MatchBuilder) AddMatches(s ...string) {
	for _, v := range s {
		mb.AddMatch(v)
	}
}

func (mb MatchBuilder) String() string {
	return mb.sb.String()
}
