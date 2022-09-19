package organize_text

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/etikett"
)

type SetKeyToEtiketten map[string]etikett.Set

func (m SetKeyToEtiketten) String() string {
	sb := &strings.Builder{}

	for h, es := range m {
		sb.WriteString(fmt.Sprintf("%s: %s\n", h, es))
	}

	return sb.String()
}

func (m SetKeyToEtiketten) Add(h string, e etikett.Etikett) {
	var es etikett.Set
	ok := false

	if es, ok = m[h]; !ok {
		es = etikett.MakeSet()
	}

	es.AddNormalized(e)
	m[h] = es
}

func (m SetKeyToEtiketten) Contains(h string, e etikett.Etikett) (ok bool) {
	var es etikett.Set

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
		for _, e := range es.Sorted() {
			out.Named.Add(z.Hinweis.String(), e)
		}
	}

	for z, _ := range a.unnamed {
		for _, e := range es.Sorted() {
			out.Unnamed.Add(z.Bezeichnung.String(), e)
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
