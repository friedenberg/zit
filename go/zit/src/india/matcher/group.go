package matcher

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/matcher_proto"
	"code.linenisgreat.com/zit/src/india/query"
)

func init() {
	gob.Register(&group{})
}

type group struct {
	konfig              schnittstellen.Konfig
	fileExtensionGetter schnittstellen.FileExtensionGetter
	expanders           kennung.Abbr

	cwd    matcher_proto.Cwd
	Hidden Matcher
	index  kennung.Index

	DefaultGattungen kennung.Gattung
	// NewQuery         *query.QueryGroup
	FDs fd.MutableSet

	dotOperatorActive bool
}

func (q group) MatcherLen() int {
	return 0
}

func (s *group) BuildQueryGroup(vs ...string) (qg matcher_proto.QueryGroup, err error) {
	var builder query.Builder

	builder.
		WithDefaultGattungen(s.DefaultGattungen).
		WithCwd(s.cwd).
		WithFileExtensionGetter(s.fileExtensionGetter).
		WithHidden(s.Hidden).
		WithExpanders(s.expanders)

	if qg, err = builder.BuildQueryGroup(vs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
