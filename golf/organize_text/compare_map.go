package organize_text

import (
	"github.com/friedenberg/zit/charlie/etikett"
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

func (in *organizeText) ToCompareMap() (out CompareMap) {
	out = CompareMap{
		Named:   make(SetEtikettenKeys),
		Unnamed: make(SetEtikettenKeys),
	}

	in.assignment.addToCompareMap(etikett.NewSet(), &out)

	return
}

func (a *assignment) addToCompareMap(es *etikett.Set, out *CompareMap) {
	es = es.Copy()
	es.Merge(a.etiketten)

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
		c.addToCompareMap(es, out)
	}

	return
}
