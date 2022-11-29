package zettel_checked_out

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/india/zettel"
	"github.com/friedenberg/zit/src/india/zettel_external"
)

type Matches struct {
	Akten, Bezeichnungen, Zettelen zettel.MutableSet
}

func MakeMatches() Matches {
	return Matches{
		Akten:         zettel.MakeMutableSetUnique(0),
		Bezeichnungen: zettel.MakeMutableSetUnique(0),
		Zettelen:      zettel.MakeMutableSetUnique(0),
	}
}

func (m Matches) appendToStringBuilder(sb *strings.Builder, ex zettel_external.Zettel) {
	typToCollection := map[gattung.Gattung]*zettel.MutableSet{
		gattung.Akte:        &m.Akten,
		gattung.Bezeichnung: &m.Bezeichnungen,
		gattung.Zettel:      &m.Zettelen,
	}

	for t, c := range typToCollection {
		if c.Len() == 1 && c.Any().Named.Stored.Objekte.Equals(&ex.Named.Stored.Objekte) {
		} else if c.Len() > 1 {
			c.Each(
				func(tz *zettel.Transacted) (err error) {
					sb.WriteString(fmt.Sprintf("\n\t%s (%s match)", tz.Named, t))
					return
				},
			)
		}
	}
}
