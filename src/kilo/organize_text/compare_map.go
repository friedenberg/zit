package organize_text

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

// TODO-P4 make generic
type SetKeyToMetadatei map[string]metadatei.Metadatei

func (m SetKeyToMetadatei) String() string {
	sb := &strings.Builder{}

	for h, es := range m {
		sb.WriteString(fmt.Sprintf("%s: %s\n", h, es))
	}

	return sb.String()
}

func (s SetKeyToMetadatei) Add(h string) {
	var m metadatei.Metadatei
	ok := false

	if m, ok = s[h]; !ok {
		m.Reset()
	}

	s[h] = m
}

func (s SetKeyToMetadatei) AddEtikett(h string, e kennung.Etikett) {
	var m metadatei.Metadatei
	ok := false

	if m, ok = s[h]; !ok {
		m.Reset()
	}

	kennung.AddNormalized(m.GetEtikettenMutable(), &e)
	s[h] = m
}

func (s SetKeyToMetadatei) ContainsEtikett(
	h string,
	e kennung.Etikett,
) (ok bool) {
	var m metadatei.Metadatei

	if m, ok = s[h]; !ok {
		return
	}

	ok = m.GetEtiketten().Contains(e)

	return
}

type CompareMap struct {
	// etikett to hinweis
	Named SetKeyToMetadatei
	// etikett to bezeichnung
	Unnamed SetKeyToMetadatei
}

func (in *Text) ToCompareMap() (out CompareMap, err error) {
	out = CompareMap{
		Named:   make(SetKeyToMetadatei),
		Unnamed: make(SetKeyToMetadatei),
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
	mes := es.CloneMutableSetPtrLike()

	var es1 kennung.EtikettSet

	if es1, err = a.expandedEtiketten(); err != nil {
		err = errors.Wrap(err)
		return
	}

	es1.Each(mes.Add)
	es = mes.CloneSetPtrLike()

	a.named.Each(
		func(z obj) (err error) {
			if z.Kennung.String() == "" {
				panic(fmt.Sprintf("%s: Kennung is nil", z))
			}

			fk := kennung.FormattedString(z.Kennung)
			out.Named.Add(fk)

			for _, e := range iter.SortedValues[kennung.Etikett](es) {
				out.Named.AddEtikett(fk, e)
			}

			for _, e := range iter.Elements[kennung.Etikett](m.EtikettSet) {
				errors.TodoP4("add typ")
				out.Named.AddEtikett(fk, e)
			}

			return
		},
	)

	a.unnamed.Each(
		func(z obj) (err error) {
			out.Unnamed.Add(z.Bezeichnung.String())

			for _, e := range iter.SortedValues[kennung.Etikett](es) {
				out.Unnamed.AddEtikett(z.Bezeichnung.String(), e)
			}

			for _, e := range iter.Elements[kennung.Etikett](m.EtikettSet) {
				errors.TodoP4("add typ")
				out.Unnamed.AddEtikett(z.Bezeichnung.String(), e)
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
