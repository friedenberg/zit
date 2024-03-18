package matcher_proto

import (
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

type Cwd interface {
	Matcher
	GetCwdFDs() fd.Set
	GetKennungForFD(*fd.FD) (*kennung.Kennung2, error)
}

type matcherCwdNop struct {
	Matcher
}

func (matcherCwdNop) GetCwdFDs() fd.Set {
	return fd.MakeSet()
}

func (matcherCwdNop) GetKennungForFD(_ *fd.FD) (*kennung.Kennung2, error) {
	return nil, nil
}

func MakeMatcherCwdNop(m Matcher) Cwd {
	return matcherCwdNop{Matcher: m}
}
