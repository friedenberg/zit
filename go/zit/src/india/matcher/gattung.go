package matcher

import (
	"fmt"

	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type Query struct {
	kennung.Sigil
	gattung.Gattung
	Matcher MatcherExactlyThisOrAllOfThese
}

func (s Query) String() string {
	return fmt.Sprintf("%s%s", s.Matcher, s.Sigil)
}

func (s Query) ContainsMatchable(m *sku.Transacted) bool {
	return s.Matcher.ContainsMatchable(m)
}

func (s Query) GetSigil() kennung.Sigil {
	return s.Sigil
}

func (s Query) GetGattung() gattung.Gattung {
	return s.Gattung
}
