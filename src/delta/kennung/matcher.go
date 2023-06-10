package kennung

import (
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/iter"
)

const (
	QueryOrOperator         = ", "
	QueryAndOperator        = " "
	QueryGroupOpenOperator  = "["
	QueryGroupCloseOperator = "]"
	QueryNegationOperator   = '^'
	QueryExactOperator      = '='
)

func init() {
	gob.Register(&matcherAnd{})
	gob.Register(&matcherOr{})
	gob.Register(&matcherNegate{})
	gob.Register(&matcherNever{})
	gob.Register(&matcherAlways{})
	gob.Register(&matcherContainsExactly{})
	gob.Register(&matcherContains{})
	gob.Register(&matcherImplicit{})
	gob.Register(&matcherExactlyThisOrAllOfThese{})
}

type Matcher interface {
	ContainsMatchable(Matchable) bool
	schnittstellen.Stringer
	MatcherLen() int
	Each(schnittstellen.FuncIter[Matcher]) error
}

type MatcherSigil interface {
	Matcher
	GetSigil() Sigil
}

type MatcherSigilPtr interface {
	MatcherSigil
	AddSigil(Sigil)
}

type MatcherKennungSansGattungWrapper interface {
	Matcher
	GetKennung() KennungSansGattung
}

type MatcherExact interface {
	Matcher
	ContainsMatchableExactly(Matchable) bool
}

type MatcherImplicit interface {
	Matcher
	GetImplicitMatcher() matcherImplicit
}

type MatcherParentPtr interface {
	Matcher
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

func VisitAllMatcherKennungSansGattungWrappers(
	f schnittstellen.FuncIter[MatcherKennungSansGattungWrapper],
	matchers ...Matcher,
) (err error) {
	return VisitAllMatchers(
		func(m Matcher) (err error) {
			if _, ok := m.(MatcherImplicit); ok {
				return iter.MakeErrStopIteration()
			}

			if mksgw, ok := m.(MatcherKennungSansGattungWrapper); ok {
				return f(mksgw)
			}

			return
		},
		matchers...,
	)
}

func VisitAllMatchers(
	f schnittstellen.FuncIter[Matcher],
	matchers ...Matcher,
) (err error) {
	for _, m := range matchers {
		if err = f(m); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		if err = m.Each(
			func(m Matcher) (err error) {
				return VisitAllMatchers(f, m)
			},
		); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

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

func (_ matcherAlways) MatcherLen() int {
	return 0
}

func (_ matcherAlways) Each(_ schnittstellen.FuncIter[Matcher]) error {
	return nil
}

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

func (_ matcherNever) MatcherLen() int {
	return 0
}

func (_ matcherNever) Each(_ schnittstellen.FuncIter[Matcher]) error {
	return nil
}

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

func (matcher matcherAnd) MatcherLen() int {
	return len(matcher.Children)
}

func (matcher *matcherAnd) Add(m Matcher) (err error) {
	matcher.Children = append(matcher.Children, m)
	return
}

func (matcher matcherAnd) String() string {
	sb := &strings.Builder{}
	sb.WriteString(QueryGroupOpenOperator)

	for i, m := range matcher.Children {
		if i > 0 {
			sb.WriteString(QueryAndOperator)
		}

		sb.WriteString(m.String())
	}

	if matcher.MatcherLen() == 0 {
		sb.WriteString(fmt.Sprintf("%t", matcher.MatchOnEmpty))
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

func (matcher matcherOr) MatcherLen() int {
	return len(matcher.Children)
}

func (matcher *matcherOr) Add(m Matcher) (err error) {
	matcher.Children = append(matcher.Children, m)
	return
}

func (matcher matcherOr) String() (out string) {
	sb := &strings.Builder{}
	sb.WriteString(QueryGroupOpenOperator)

	for i, m := range matcher.Children {
		if i > 0 {
			sb.WriteString(QueryOrOperator)
		}

		sb.WriteString(m.String())
	}

	if matcher.MatcherLen() == 0 {
		sb.WriteString(fmt.Sprintf("%t", matcher.MatchOnEmpty))
	}

	sb.WriteString(QueryGroupCloseOperator)

	out = sb.String()
	return
}

func (matcher matcherOr) ContainsMatchable(matchable Matchable) bool {
	if len(matcher.Children) == 0 {
		return matcher.MatchOnEmpty
	}

	l := 0

	for _, m := range matcher.Children {
		if m.ContainsMatchable(matchable) {
			return true
		}

		l += m.MatcherLen()
	}

	if l == 0 && matcher.MatchOnEmpty {
		return true
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

//    ____            _        _
//   / ___|___  _ __ | |_ __ _(_)_ __  ___
//  | |   / _ \| '_ \| __/ _` | | '_ \/ __|
//  | |__| (_) | | | | || (_| | | | | \__ \
//   \____\___/|_| |_|\__\__,_|_|_| |_|___/
//

func MakeMatcherContains(
	k KennungSansGattung,
	ki Index,
) MatcherKennungSansGattungWrapper {
	return &matcherContains{Kennung: k, index: ki}
}

type matcherContains struct {
	Kennung KennungSansGattung
	index   Index
}

func (matcher matcherContains) GetKennung() KennungSansGattung {
	return matcher.Kennung
}

func (m matcherContains) MatcherLen() int {
	return 0
}

func (_ matcherContains) Each(_ schnittstellen.FuncIter[Matcher]) error {
	return nil
}

func (matcher matcherContains) String() string {
	if matcher.Kennung == nil {
		return ""
	}

	return FormattedString(matcher.Kennung)
}

func (matcher matcherContains) ContainsMatchable(matchable Matchable) bool {
	if !KennungContainsMatchable(matcher.Kennung, matchable, matcher.index) {
		return false
	}

	return true
}

//    ____            _        _           _____                _   _
//   / ___|___  _ __ | |_ __ _(_)_ __  ___| ____|_  ____ _  ___| |_| |_   _
//  | |   / _ \| '_ \| __/ _` | | '_ \/ __|  _| \ \/ / _` |/ __| __| | | | |
//  | |__| (_) | | | | || (_| | | | | \__ \ |___ >  < (_| | (__| |_| | |_| |
//   \____\___/|_| |_|\__\__,_|_|_| |_|___/_____/_/\_\__,_|\___|\__|_|\__, |
//                                                                    |___/

func MakeMatcherContainsExactly(
	k KennungSansGattung,
) MatcherKennungSansGattungWrapper {
	return &matcherContainsExactly{Kennung: k}
}

type matcherContainsExactly struct {
	Kennung KennungSansGattung
}

func (matcher matcherContainsExactly) GetKennung() KennungSansGattung {
	return matcher.Kennung
}

func (m matcherContainsExactly) MatcherLen() int {
	return 0
}

func (_ matcherContainsExactly) Each(_ schnittstellen.FuncIter[Matcher]) error {
	return nil
}

func (matcher matcherContainsExactly) String() string {
	if matcher.Kennung == nil {
		return ""
	}

	return FormattedString(matcher.Kennung) + string(QueryExactOperator)
}

func (matcher matcherContainsExactly) ContainsMatchable(
	matchable Matchable,
) bool {
	return KennungContainsExactlyMatchable(matcher.Kennung, matchable)
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

func (matcher matcherNegate) MatcherLen() int {
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
	ok := !matcher.Child.ContainsMatchable(matchable)

	return ok
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

func (matcher matcherImplicit) GetImplicitMatcher() matcherImplicit {
	return matcher
}

func (matcher matcherImplicit) MatcherLen() int {
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
	return matcher.Child.String()
	// return ""
}

func (matcher matcherImplicit) ContainsMatchable(matchable Matchable) bool {
	return matcher.Child.ContainsMatchable(matchable)
}

func (matcher matcherImplicit) Each(f schnittstellen.FuncIter[Matcher]) error {
	return f(matcher.Child)
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

func (m matcherGattung) MatcherLen() int {
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
	if m.MatcherLen() == 0 {
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

//  __        ___ _   _     ____  _       _ _
//  \ \      / (_) |_| |__ / ___|(_) __ _(_) |
//   \ \ /\ / /| | __| '_ \\___ \| |/ _` | | |
//    \ V  V / | | |_| | | |___) | | (_| | | |
//     \_/\_/  |_|\__|_| |_|____/|_|\__, |_|_|
//                                  |___/

func MakeMatcherWithSigil(m Matcher, s Sigil) MatcherSigilPtr {
	return &matcherWithSigil{
		Matcher: m,
		Sigil:   s,
	}
}

type matcherWithSigil struct {
	Sigil
	Matcher
}

func (m matcherWithSigil) Len() int {
	if m.Matcher == nil {
		return 0
	}

	return 1
}

func (m matcherWithSigil) String() string {
	sb := &strings.Builder{}

	if m.Matcher != nil {
		sb.WriteString(m.Matcher.String())
	}

	sb.WriteString(m.Sigil.String())

	return sb.String()
}

func (m matcherWithSigil) GetSigil() Sigil {
	return m.Sigil
}

func (m *matcherWithSigil) AddSigil(v Sigil) {
	errors.TodoP1("add sigils to children")
	m.Sigil.Add(v)
}

func (m *matcherWithSigil) Add(child Matcher) (err error) {
	m.Matcher = child
	return
}

func (matcher matcherWithSigil) ContainsMatchable(matchable Matchable) bool {
	if matcher.Matcher == nil {
		return true
	}

	return matcher.Matcher.ContainsMatchable(matchable)
}

func (matcher matcherWithSigil) Each(f schnittstellen.FuncIter[Matcher]) error {
	return f(matcher.Matcher)
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
