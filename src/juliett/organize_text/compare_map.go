package organize_text

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/etikett"
)

type TupleEtikettKey struct {
	Etikett, Key string
}

type SetEtikettenKeys map[TupleEtikettKey]bool

func (m SetEtikettenKeys) Add(e, k string) {
	t := TupleEtikettKey{
		Etikett: e,
		Key:     k,
	}

	m[t] = true
}

type CompareMap struct {
	// etikett to hinweis
	Named SetEtikettenKeys
	// etikett to bezeichnung
	Unnamed SetEtikettenKeys
}

func (in *organizeText) ToCompareMap() (out CompareMap, err error) {
	out = CompareMap{
		Named:   make(SetEtikettenKeys),
		Unnamed: make(SetEtikettenKeys),
	}

	if err = in.assignment.addToCompareMap(etikett.NewSet(), &out); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *assignment) addToCompareMap(es *etikett.Set, out *CompareMap) (err error) {
	es = es.Copy()

	var es1 etikett.Set

	if es1, err = a.expandedEtiketten(); err != nil {
		err = errors.Wrap(err)
		return
	}

	es.Merge(es1)
	errors.Print(es)

	for z, _ := range a.named {
		for e, _ := range *es {
			out.Named.Add(e, z.Hinweis)
		}
	}

	for z, _ := range a.unnamed {
		for e, _ := range *es {
			out.Unnamed.Add(e, z.Bezeichnung)
		}
	}

	for _, c := range a.children {
		if err = c.addToCompareMap(es, out); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
