package etiketten_path

import (
	"slices"

	"code.linenisgreat.com/zit/src/alfa/errors"
)

type Etiketten struct {
	Paths []*Path
	All   []EtikettWithParents
}

func (es *Etiketten) Reset() {
	// TODO pool *Path's
	es.Paths = es.Paths[:0]
	es.All = es.All[:0]
}

func (a *Etiketten) ResetWith(b *Etiketten) {
	a.Reset()
	a.Paths = slices.Grow(a.Paths, b.Len())
	copy(a.Paths, b.Paths)

	a.Reset()
	a.All = slices.Grow(a.All, b.Len())
	copy(a.All, b.All)
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
