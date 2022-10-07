package organize_text

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/etikett"
)

type SetKeyToEtiketten map[string]etikett.MutableSet

func (m SetKeyToEtiketten) String() string {
	sb := &strings.Builder{}

	for h, es := range m {
		sb.WriteString(fmt.Sprintf("%s: %s\n", h, es))
	}

	return sb.String()
}

func (m SetKeyToEtiketten) Add(h string, e etikett.Etikett) {
	var es etikett.MutableSet
	ok := false

	if es, ok = m[h]; !ok {
		es = etikett.MakeMutableSet()
	}

	es.AddNormalized(e)
	m[h] = es
}

func (m SetKeyToEtiketten) Contains(h string, e etikett.Etikett) (ok bool) {
	var es etikett.MutableSet

	if es, ok = m[h]; !ok {
		return
	}

	ok = es.Contains(e)

	return
}

type CompareMap struct {
	// etikett to hinweis
	Named SetKeyToEtiketten
	// etikett to bezeichnung
	Unnamed SetKeyToEtiketten
}

func (in *Text) ToCompareMap() (out CompareMap, err error) {
	out = CompareMap{
		Named:   make(SetKeyToEtiketten),
		Unnamed: make(SetKeyToEtiketten),
	}

	if err = in.assignment.addToCompareMap(in.Metadatei, etikett.MakeSet(), &out); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *assignment) addToCompareMap(m Metadatei, es etikett.Set, out *CompareMap) (err error) {
	mes := es.MutableCopy()

	var es1 etikett.Set

	if es1, err = a.expandedEtiketten(); err != nil {
		err = errors.Wrap(err)
		return
	}

	mes.Merge(es1)
	es = mes.Copy()

	for z, _ := range a.named {
		for _, e := range es.Sorted() {
			out.Named.Add(z.Hinweis.String(), e)
		}

		for _, e := range m.Set.Etiketten() {
			//TODO add typ
			out.Named.Add(z.Hinweis.String(), e)
		}
	}

	for z, _ := range a.unnamed {
		for _, e := range es.Sorted() {
			out.Unnamed.Add(z.Bezeichnung.String(), e)
		}

		for _, e := range m.Set.Etiketten() {
			//TODO add typ
			out.Unnamed.Add(z.Bezeichnung.String(), e)
		}
	}

	for _, c := range a.children {
		if err = c.addToCompareMap(m, es, out); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
