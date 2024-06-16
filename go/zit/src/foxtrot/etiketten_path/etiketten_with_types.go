package etiketten_path

import (
	"fmt"
	"slices"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

type Etiketten struct {
	Paths PathsWithTypes // TODO implement
	All   EtikettenWithParentsAndTypes
}

func (a *Etiketten) String() string {
	return fmt.Sprintf("[Paths: %s, All: %s]", a.Paths, a.All)
}

func (a *Etiketten) Reset() {
	// TODO pool *Path's
	a.Paths.Reset()
	a.All.Reset()
}

// TODO improve performance
func (a *Etiketten) ResetWith(b *Etiketten) {
	a.Paths = slices.Grow(a.Paths, len(b.Paths))

	for _, p := range b.Paths {
		a.AddPath(p.Clone())
	}
	// a.Paths = a.Paths[:cap(a.Paths)]
	// nPaths := copy(a.Paths, b.Paths)

	// a.All = slices.Grow(a.All, len(b.All))
	// a.All = a.All[:cap(a.All)]
	// nAll := copy(a.All, b.All)
	// ui.Debug().Print(nPaths, nAll, a, b)
}

func (a *Etiketten) AddSuperFrom(
	b *Etiketten,
	prefix *Etikett,
) (err error) {
	for _, ep := range b.Paths {
		ui.Log().Print("adding", prefix, ep)
		if prefix.ComparePartial(ep.First()) == 0 {
			continue
		}

		prefixPath := makePath(prefix)
		prefixPath.Add(ep.Path...)

		c := &PathWithType{
			Path: prefixPath,
			Type: TypeSuper,
		}

		if err = a.AddPath(c); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (es *Etiketten) AddEtikettOld(e kennung.Etikett) (err error) {
	return es.AddEtikett(catgut.MakeFromString(e.String()))
}

func (es *Etiketten) AddEtikett(e *Etikett) (err error) {
	if e.IsEmpty() {
		return
	}

	p := MakePathWithType(e)

	if err = es.AddPathWithType(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es *Etiketten) AddSelf(e *Etikett) (err error) {
	if e.IsEmpty() {
		return
	}

	p := MakePathWithType(e)
  p.Type = TypeSelf

	if err = es.AddPathWithType(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es *Etiketten) AddPathWithType(pwt *PathWithType) (err error) {
	_, alreadyExists := es.Paths.AddPath(pwt)

	if alreadyExists {
		return
	}

	for _, e := range pwt.Path {
		if err = es.All.Add(e, pwt); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (es *Etiketten) AddPath(p *PathWithType) (err error) {
	if err = es.AddPathWithType(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Etiketten) Set(v string) (err error) {
	vs := strings.Split(v, ",")

	for _, v := range vs {
		var e kennung.Etikett

		if err = e.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}

		es := catgut.MakeFromString(e.String())

		if err = s.AddEtikett(es); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
