package etiketten_path

import (
	"slices"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

type Etiketten struct {
	Paths []*Path
	All   []EtikettWithParents
}

func (a *Etiketten) Reset() {
	// TODO pool *Path's
	a.Paths = a.Paths[:0]
	a.All = a.All[:0]
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

func (es *Etiketten) Len() int {
	return len(es.Paths)
}

func (es *Etiketten) Less(i, j int) bool {
	return es.Paths[i].Compare(es.Paths[j]) == -1
}

func (es *Etiketten) Swap(i, j int) {
	es.Paths[j], es.Paths[i] = es.Paths[i], es.Paths[j]
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
	if p.IsEmpty() {
		return
	}

	idx, ok := es.ContainsPath(p)

	if ok {
		return
	}

	es.Paths = slices.Insert(es.Paths, idx, p)

	for _, e := range *p {
		if err = es.addToAll(e, p); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (es *Etiketten) addToAll(e *Etikett, p *Path) (err error) {
	idx, ok := es.ContainsEtikett(e)

	var a EtikettWithParents

	if ok {
		a = es.All[idx]
		a.AddParent(p)
		es.All[idx] = a
	} else {
		a = EtikettWithParents{Etikett: e}
		a.AddParent(p)

		if idx == len(es.All) {
			es.All = append(es.All, a)
		} else {
			es.All = slices.Insert(es.All, idx, a)
		}
	}

	return
}

func (es *Etiketten) ContainsPath(p *Path) (int, bool) {
	return slices.BinarySearchFunc(
		es.Paths,
		p,
		func(ep *Path, el *Path) int {
			return ep.Compare(p)
		},
	)
}

func (es *Etiketten) GetMatching(e *Etikett) (matching []EtikettWithParents) {
	i, ok := es.ContainsEtikett(e)

	if !ok {
		return
	}

	for _, ewp := range es.All[i:] {
		cmp := ewp.ComparePartial(e)

		if cmp != 0 {
			return
		}

		matching = append(matching, ewp)
	}

	return
}

func (es *Etiketten) ContainsEtikett(e *Etikett) (i int, ok bool) {
	i, ok = slices.BinarySearchFunc(
		es.All,
		e,
		func(ewp EtikettWithParents, e *Etikett) int {
			cmp := ewp.Etikett.ComparePartial(e)
			return cmp
		},
	)

	return
}
