package matcher

import (
	"encoding/gob"
	"strings"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

func init() {
	gob.Register(&matcherSigil{})
}

//   __  __       _       _                 _   _ _     _     _
//  |  \/  | __ _| |_ ___| |__   ___ _ __  | | | (_) __| | __| | ___ _ __
//  | |\/| |/ _` | __/ __| '_ \ / _ \ '__| | |_| | |/ _` |/ _` |/ _ \ '_ \
//  | |  | | (_| | || (__| | | |  __/ |    |  _  | | (_| | (_| |  __/ | | |
//  |_|  |_|\__,_|\__\___|_| |_|\___|_|    |_| |_|_|\__,_|\__,_|\___|_| |_|
//

func MakeMatcherSigil(s kennung.Sigil, m Matcher) MatcherSigilPtr {
	return &matcherSigil{
		MatchSigil: s,
		Matcher:    m,
	}
}

func MakeMatcherSigilMatchOnMissing(
	s kennung.Sigil,
	m Matcher,
) MatcherSigilPtr {
	return &matcherSigil{
		MatchSigil:     s,
		Matcher:        m,
		MatchOnMissing: true,
	}
}

type matcherSigil struct {
	MatchSigil kennung.Sigil
	kennung.Sigil
	Matcher
	MatchOnMissing bool
}

func (m matcherSigil) Len() int {
	if m.Matcher == nil {
		return 0
	}

	return 1
}

func (m matcherSigil) String() string {
	sb := &strings.Builder{}

	if m.Matcher != nil {
		sb.WriteString(m.Matcher.String())
	}

	sb.WriteString(m.Sigil.String())

	return sb.String()
}

func (m matcherSigil) GetSigil() kennung.Sigil {
	return m.Sigil
}

func (m *matcherSigil) AddSigil(v kennung.Sigil) {
	m.Sigil.Add(v)
}

func (m *matcherSigil) Add(child Matcher) (err error) {
	m.Matcher = child
	return
}

func (matcher matcherSigil) ContainsMatchable(matchable *sku.Transacted) bool {
	switch {
	case matcher.MatchOnMissing && !matcher.Contains(matcher.MatchSigil):
		fallthrough

	case matcher.Contains(matcher.MatchSigil):
		fallthrough

	case matcher.Matcher == nil:
		return true
	}

	return matcher.Matcher.ContainsMatchable(matchable)
}

func (matcher matcherSigil) Each(f schnittstellen.FuncIter[Matcher]) error {
	return f(matcher.Matcher)
}

//   __  __       _       _                 _   _ _     _     _
//  |  \/  | __ _| |_ ___| |__   ___ _ __  | | | (_) __| | __| | ___ _ __
//  | |\/| |/ _` | __/ __| '_ \ / _ \ '__| | |_| | |/ _` |/ _` |/ _ \ '_ \
//  | |  | | (_| | || (__| | | |  __/ |    |  _  | | (_| | (_| |  __/ | | |
//  |_|  |_|\__,_|\__\___|_| |_|\___|_|    |_| |_|_|\__,_|\__,_|\___|_| |_|
//

func MakeMatcherExcludeHidden(m Matcher, s kennung.Sigil) MatcherSigilPtr {
	return &matcherExcludeHidden{
		Sigil:   s,
		Matcher: m,
	}
}

type matcherExcludeHidden struct {
	Sigil   kennung.Sigil
	Matcher Matcher
}

func (m matcherExcludeHidden) MatcherLen() int {
	if m.Matcher == nil {
		return 0
	}

	return 1
}

func (m matcherExcludeHidden) String() string {
	sb := &strings.Builder{}

	if m.Matcher != nil {
		sb.WriteString(m.Matcher.String())
	}

	sb.WriteString(m.Sigil.String())

	return sb.String()
}

func (m matcherExcludeHidden) GetSigil() kennung.Sigil {
	return m.Sigil
}

func (m *matcherExcludeHidden) AddSigil(v kennung.Sigil) {
	m.Sigil.Add(v)
}

func (pred matcherExcludeHidden) ContainsMatchable(
	val *sku.Transacted,
) bool {
	if pred.Sigil.IncludesHidden() {
		return true
	}

	if pred.Matcher == nil {
		return true
	}

	if !pred.Matcher.ContainsMatchable(val) {
		return true
	}

	return false
}

func (matcher matcherExcludeHidden) Each(
	f schnittstellen.FuncIter[Matcher],
) error {
	return f(matcher.Matcher)
}
