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
	Paths SlicePaths
	All   SliceEtikettWithParents
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

func (a *Etiketten) AddFrom(b *Etiketten, prefix *Path) (err error) {
	for _, ep := range b.Paths {
		ui.Log().Print("adding", prefix, ep)
		if prefix.First().ComparePartial(ep.First()) == 0 {
			continue
		}

		if err = a.AddPath(prefix.CloneAndAddPath(ep)); err != nil {
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

	p := MakePath(e)

	if err = es.AddPath(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es *Etiketten) AddPath(p *Path) (err error) {
	_, alreadyExists := es.Paths.AddPath(p)

	if alreadyExists {
		return
	}

	for _, e := range *p {
		if err = es.All.Add(e, p); err != nil {
			err = errors.Wrap(err)
			return
		}
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
