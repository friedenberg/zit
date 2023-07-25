package organize_text

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
)

// TODO-P4 make generic
type SetKeyToEtiketten map[string]kennung.EtikettMutableSet

func (m SetKeyToEtiketten) String() string {
	sb := &strings.Builder{}

	for h, es := range m {
		sb.WriteString(fmt.Sprintf("%s: %s\n", h, es))
	}

	return sb.String()
}

func (m SetKeyToEtiketten) Add(h string, e kennung.Etikett) {
	var es kennung.EtikettMutableSet
	ok := false

	if es, ok = m[h]; !ok {
		es = kennung.MakeEtikettMutableSet()
	}

	kennung.AddNormalized(es, e)
	m[h] = es
}

func (m SetKeyToEtiketten) Contains(h string, e kennung.Etikett) (ok bool) {
	var es kennung.EtikettMutableSet

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

	if err = in.assignment.addToCompareMap(
		in.Metadatei,
		kennung.MakeEtikettSet(),
		&out,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *assignment) addToCompareMap(
	m Metadatei,
	es kennung.EtikettSet,
	out *CompareMap,
) (err error) {
	mes := es.CloneMutableSetLike()

	var es1 kennung.EtikettSet

	if es1, err = a.expandedEtiketten(); err != nil {
		err = errors.Wrap(err)
		return
	}

	es1.Each(mes.Add)
	es = mes.CloneSetLike()

	a.named.Each(
		func(z obj) (err error) {
			for _, e := range collections.SortedValues(es) {
				if z.Kennung == nil {
					panic(fmt.Sprintf("%s: Kennung is nil", z))
				}

				out.Named.Add(kennung.FormattedString(z.Kennung), e)
			}

			for _, e := range m.EtikettSet.Elements() {
				errors.TodoP4("add typ")
				out.Named.Add(kennung.FormattedString(z.Kennung), e)
			}

			return
		},
	)

	a.unnamed.Each(
		func(z obj) (err error) {
			for _, e := range collections.SortedValues(es) {
				out.Unnamed.Add(z.Bezeichnung.String(), e)
			}

			for _, e := range m.EtikettSet.Elements() {
				errors.TodoP4("add typ")
				out.Unnamed.Add(z.Bezeichnung.String(), e)
			}

			return
		},
	)

	for _, c := range a.children {
		if err = c.addToCompareMap(m, es, out); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
