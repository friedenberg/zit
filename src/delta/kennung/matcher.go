package kennung

import (
	"encoding/gob"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/iter"
)

const (
	QueryOrOperator         = ", "
	QueryAndOperator        = " "
	QueryGroupOpenOperator  = "]"
	QueryGroupCloseOperator = "["
)

func init() {
	gob.Register(&matcherAnd{})
	gob.Register(&matcherOr{})
	gob.Register(&matcherNegate{})
	gob.Register(&matcherNever{})
	gob.Register(&matcherAlways{})
	gob.Register(&matcherImpExp{})
}

type Matcher interface {
	ContainsMatchable(Matchable) bool
	String() string
	// schnittstellen.Stringer
}

type MatcherParent interface {
	Matcher
	Len() int
	Each(schnittstellen.FuncIter[Matcher]) error
}

type MatcherParentPtr interface {
	MatcherParent
	Add(Matcher) error
}

func LenMatchers(
	matchers ...Matcher,
) (i int) {
	inc := func(m Matcher) (err error) {
		if _, ok := m.(Kennung); ok {
			i++
		}

		return
	}

	VisitAllMatchers(inc, matchers...)

	return
}

func VisitAllMatchers(
	f schnittstellen.FuncIter[Matcher],
	matchers ...Matcher,
) (err error) {
	for _, m := range matchers {
		if err = f(m); err != nil {
			err = errors.Wrap(err)
			return
		}

		mp, ok := m.(MatcherParent)

		if !ok {
			continue
		}

		if err = mp.Each(
			func(m Matcher) (err error) {
				return VisitAllMatchers(f, m)
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

//      _    _
//     / \  | |_      ____ _ _   _ ___
//    / _ \ | \ \ /\ / / _` | | | / __|
//   / ___ \| |\ V  V / (_| | |_| \__ \
//  /_/   \_\_| \_/\_/ \__,_|\__, |___/
//                           |___/

func MakeMatcherAlways() Matcher {
	return matcherAlways{}
}

type matcherAlways struct{}

func (_ matcherAlways) String() string {
	return "ALWAYS"
}

func (_ matcherAlways) ContainsMatchable(_ Matchable) bool {
	return true
}

//   _   _
//  | \ | | _____   _____ _ __
//  |  \| |/ _ \ \ / / _ \ '__|
//  | |\  |  __/\ V /  __/ |
//  |_| \_|\___| \_/ \___|_|
//

func MakeMatcherNever() Matcher {
	return matcherNever{}
}

type matcherNever struct{}

func (_ matcherNever) String() string {
	return "NEVER"
}

func (_ matcherNever) ContainsMatchable(_ Matchable) bool {
	return false
}

//      _              _
//     / \   _ __   __| |
//    / _ \ | '_ \ / _` |
//   / ___ \| | | | (_| |
//  /_/   \_\_| |_|\__,_|
//

func MakeMatcherAnd(ms ...Matcher) MatcherParentPtr {
	return &matcherAnd{
		MatchOnEmpty: true,
		Children:     ms,
	}
}

func MakeMatcherAndDoNotMatchOnEmpty(ms ...Matcher) MatcherParentPtr {
	return &matcherAnd{
		Children: ms,
	}
}

type matcherAnd struct {
	MatchOnEmpty bool
	Children     []Matcher
}

func (matcher matcherAnd) Len() int {
	return len(matcher.Children)
}

func (matcher *matcherAnd) Add(m Matcher) (err error) {
	matcher.Children = append(matcher.Children, m)
	return
}

func (matcher matcherAnd) String() string {
	if matcher.Len() == 0 {
		return ""
	}

	sb := &strings.Builder{}
	sb.WriteString(QueryGroupOpenOperator)

	for i, m := range matcher.Children {
		if i > 0 {
			sb.WriteString(QueryAndOperator)
		}

		sb.WriteString(m.String())
	}

	sb.WriteString(QueryGroupCloseOperator)

	return sb.String()
}

func (matcher matcherAnd) ContainsMatchable(matchable Matchable) bool {
	if len(matcher.Children) == 0 {
		return matcher.MatchOnEmpty
	}

	for _, m := range matcher.Children {
		if !m.ContainsMatchable(matchable) {
			return false
		}
	}

	return true
}

func (matcher matcherAnd) Each(f schnittstellen.FuncIter[Matcher]) (err error) {
	for _, m := range matcher.Children {
		if err = f(m); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

//    ___
//   / _ \ _ __
//  | | | | '__|
//  | |_| | |
//   \___/|_|
//

func MakeMatcherOr(ms ...Matcher) MatcherParentPtr {
	return &matcherOr{
		MatchOnEmpty: true,
		Children:     ms,
	}
}

func MakeMatcherOrDoNotMatchOnEmpty(ms ...Matcher) MatcherParentPtr {
	return &matcherOr{
		Children: ms,
	}
}

type matcherOr struct {
	MatchOnEmpty bool
	Children     []Matcher
}

func (matcher matcherOr) Len() int {
	return len(matcher.Children)
}

func (matcher *matcherOr) Add(m Matcher) (err error) {
	matcher.Children = append(matcher.Children, m)
	return
}

func (matcher matcherOr) String() (out string) {
	if matcher.Len() == 0 {
		return
	}

	sb := &strings.Builder{}
	sb.WriteString(QueryGroupOpenOperator)

	for i, m := range matcher.Children {
		if i > 0 {
			sb.WriteString(QueryOrOperator)
		}

		sb.WriteString(m.String())
	}

	sb.WriteString(QueryGroupCloseOperator)

	out = sb.String()
	return
}

func (matcher matcherOr) ContainsMatchable(matchable Matchable) bool {
	if len(matcher.Children) == 0 {
		return matcher.MatchOnEmpty
	}

	for _, m := range matcher.Children {
		if m.ContainsMatchable(matchable) {
			return true
		}
	}

	return false
}

func (matcher matcherOr) Each(f schnittstellen.FuncIter[Matcher]) (err error) {
	for _, m := range matcher.Children {
		if err = f(m); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

//   _   _                  _
//  | \ | | ___  __ _  __ _| |_ ___
//  |  \| |/ _ \/ _` |/ _` | __/ _ \
//  | |\  |  __/ (_| | (_| | ||  __/
//  |_| \_|\___|\__, |\__,_|\__\___|
//              |___/

func MakeMatcherNegate(m Matcher) MatcherParentPtr {
	return &matcherNegate{Child: m}
}

type matcherNegate struct {
	Child Matcher
}

func (matcher matcherNegate) Len() int {
	if matcher.Child == nil {
		return 0
	}

	return 1
}

func (matcher *matcherNegate) Add(m Matcher) error {
	matcher.Child = m
	return nil
}

func (matcher matcherNegate) String() string {
	if matcher.Child == nil {
		return ""
	}

	return string(QueryNegationOperator) + matcher.Child.String()
}

func (matcher matcherNegate) ContainsMatchable(matchable Matchable) bool {
	return !matcher.Child.ContainsMatchable(matchable)
}

func (matcher matcherNegate) Each(f schnittstellen.FuncIter[Matcher]) error {
	return f(matcher.Child)
}

//   ___                 _ _      _ _
//  |_ _|_ __ ___  _ __ | (_) ___(_) |_
//   | || '_ ` _ \| '_ \| | |/ __| | __|
//   | || | | | | | |_) | | | (__| | |_
//  |___|_| |_| |_| .__/|_|_|\___|_|\__|
//                |_|

func MakeMatcherImplicit(m Matcher) MatcherParentPtr {
	return &matcherImplicit{Child: m}
}

type matcherImplicit struct {
	Child Matcher
}

func (matcher matcherImplicit) Len() int {
	if matcher.Child == nil {
		return 0
	}

	return 1
}

func (matcher *matcherImplicit) Add(m Matcher) error {
	matcher.Child = m
	return nil
}

func (matcher matcherImplicit) String() string {
	if matcher.Child == nil {
		return ""
	}

	return matcher.Child.String()
}

func (matcher matcherImplicit) ContainsMatchable(matchable Matchable) bool {
	return matcher.Child.ContainsMatchable(matchable)
}

func (matcher matcherImplicit) Each(f schnittstellen.FuncIter[Matcher]) error {
	return nil
	// return f(matcher.Child)
}

//    ____       _   _
//   / ___| __ _| |_| |_ _   _ _ __   __ _
//  | |  _ / _` | __| __| | | | '_ \ / _` |
//  | |_| | (_| | |_| |_| |_| | | | | (_| |
//   \____|\__,_|\__|\__|\__,_|_| |_|\__, |
//                                   |___/

func MakeMatcherGattung(m map[gattung.Gattung]Matcher) *matcherGattung {
	if m == nil {
		m = make(map[gattung.Gattung]Matcher)
	}

	return &matcherGattung{Children: m}
}

type matcherGattung struct {
	Children map[gattung.Gattung]Matcher
}

func (m matcherGattung) Len() int {
	return len(m.Children)
}

func (m *matcherGattung) Set(g gattung.Gattung, child Matcher) error {
	c1, ok := m.Children[g]

	if ok && c1 != nil {
		c1 = MakeMatcherAnd(c1, child)
	} else {
		c1 = child
	}

	m.Children[g] = c1

	return nil
}

func (m matcherGattung) String() string {
	if m.Len() == 0 {
		return ""
	}

	sb := &strings.Builder{}
	hasAny := false

	for g, child := range m.Children {
		if hasAny == true {
			sb.WriteString(QueryAndOperator)
		}

		sb.WriteString(child.String())
		sb.WriteString(g.String())
	}

	return sb.String()
}

func (matcher matcherGattung) ContainsMatchable(matchable Matchable) bool {
	g := gattung.Make(matchable.GetGattung())

	m, ok := matcher.Children[g]

	if !ok {
		return false
	}

	return m.ContainsMatchable(matchable)
}

func (matcher matcherGattung) Each(
	f schnittstellen.FuncIter[Matcher],
) (err error) {
	for _, m := range matcher.Children {
		if err = f(m); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

//   ___                 _____
//  |_ _|_ __ ___  _ __ | ____|_  ___ __
//   | || '_ ` _ \| '_ \|  _| \ \/ / '_ \
//   | || | | | | | |_) | |___ >  <| |_) |
//  |___|_| |_| |_| .__/|_____/_/\_\ .__/
//                |_|              |_|

func MakeMatcherImpExp(
	imp Matcher,
	exp MatcherParentPtr,
) *matcherImpExp {
	return &matcherImpExp{
		Implicit: imp,
		Explicit: exp,
	}
}

type matcherImpExp struct {
	Implicit Matcher
	Explicit MatcherParentPtr
}

func (m matcherImpExp) Len() (i int) {
	if m.Explicit != nil && m.Explicit.Len() > 0 {
		i++
	}

	return
}

func (m *matcherImpExp) Add(child Matcher) error {
	return m.Explicit.Add(child)
}

func (m matcherImpExp) String() string {
	if m.Explicit == nil {
		return ""
	}

	return m.Explicit.String()
}

func (matcher matcherImpExp) ContainsMatchable(matchable Matchable) bool {
	if matcher.Implicit != nil && !matcher.Implicit.ContainsMatchable(matchable) {
		return false
	}

	if matcher.Explicit != nil && !matcher.Explicit.ContainsMatchable(matchable) {
		return false
	}

	return true
}

func (matcher matcherImpExp) Each(
	f schnittstellen.FuncIter[Matcher],
) (err error) {
	// consider using a flag class like "ImplicitMatcher" to mark Imp rather than
	// breaking the rules of `Each`
	// if matcher.Implicit != nil {
	// 	if err = f(matcher.Implicit); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}
	// }

	if matcher.Explicit != nil {
		if err = f(matcher.Explicit); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func MakeMatcherFuncIter[T Matchable](m Matcher) schnittstellen.FuncIter[T] {
	return func(e T) (err error) {
		if !m.ContainsMatchable(e) {
			err = iter.MakeErrStopIteration()
			return
		}

		return
	}
}
