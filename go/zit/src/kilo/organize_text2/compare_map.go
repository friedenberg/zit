package organize_text2

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/echo/bezeichnung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
)

type SetKeyToMetadatei map[string]*metadatei.Metadatei

func (m SetKeyToMetadatei) String() string {
	sb := &strings.Builder{}

	for h, es := range m {
		fmt.Fprintf(sb, "%s: %s\n", h, es)
	}

	return sb.String()
}

func (s SetKeyToMetadatei) Add(h string, b bezeichnung.Bezeichnung) {
	var m *metadatei.Metadatei
	ok := false

	if m, ok = s[h]; !ok {
		m = &metadatei.Metadatei{}
		metadatei.Resetter.Reset(m)
		m.Bezeichnung = b
	}

	s[h] = m
}

func (s SetKeyToMetadatei) AddEtikett(
	h string,
	e kennung.Etikett,
	b bezeichnung.Bezeichnung,
) {
	var m *metadatei.Metadatei
	ok := false

	if m, ok = s[h]; !ok {
		metadatei.Resetter.Reset(m)
		m.Bezeichnung = b
	}

	if !bezeichnung.Equaler.Equals(m.Bezeichnung, b) {
		panic(fmt.Sprintf("bezeichnung changes: %q != %q", m.Bezeichnung, b))
	}

	kennung.AddNormalized(m.GetEtikettenMutable(), &e)

	s[h] = m
}

func (s SetKeyToMetadatei) ContainsEtikett(
	h string,
	e kennung.Etikett,
) (ok bool) {
	var m *metadatei.Metadatei

	if m, ok = s[h]; !ok {
		return
	}

	ok = m.GetEtiketten().Contains(e)

	return
}

type CompareMap struct {
	Named   SetKeyToMetadatei // etikett to hinweis
	Unnamed SetKeyToMetadatei // etikett to bezeichnung
}

func (in *Text) ToCompareMap() (out CompareMap, err error) {
	out = CompareMap{
		Named:   make(SetKeyToMetadatei),
		Unnamed: make(SetKeyToMetadatei),
	}

	if err = in.addToCompareMap(
		in,
		in.Metadatei,
		kennung.MakeEtikettSet(),
		&out,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Assignment) addToCompareMap(
	ot *Text,
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

	if err = a.Named.Each(
		func(z *obj) (err error) {
			if z.Kennung.String() == "" {
				panic(fmt.Sprintf("%s: Kennung is nil", z))
			}

			fk := kennung.FormattedString(&z.Kennung)
			out.Named.Add(fk, z.Metadatei.Bezeichnung)

			for _, e := range iter.SortedValues[kennung.Etikett](es) {
				out.Named.AddEtikett(fk, e, z.Metadatei.Bezeichnung)
			}

			for _, e := range iter.Elements[kennung.Etikett](m.EtikettSet) {
				errors.TodoP4("add typ")
				out.Named.AddEtikett(fk, e, z.Metadatei.Bezeichnung)
			}

			if ot.Konfig.NewOrganize {
				if err = z.Metadatei.GetEtiketten().EachPtr(
					func(e *kennung.Etikett) (err error) {
						if a.Contains(e) {
							return
						}

						out.Named.AddEtikett(fk, *e, z.Metadatei.Bezeichnung)
						return
					},
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = a.Unnamed.Each(
		func(z *obj) (err error) {
			out.Unnamed.Add(
				z.Metadatei.Bezeichnung.String(),
				z.Metadatei.Bezeichnung,
			)

			for _, e := range iter.SortedValues[kennung.Etikett](es) {
				out.Unnamed.AddEtikett(
					z.Metadatei.Bezeichnung.String(),
					e,
					z.Metadatei.Bezeichnung,
				)
			}

			for _, e := range iter.Elements[kennung.Etikett](m.EtikettSet) {
				errors.TodoP4("add typ")
				out.Unnamed.AddEtikett(
					z.Metadatei.Bezeichnung.String(),
					e,
					z.Metadatei.Bezeichnung,
				)
			}

			if ot.Konfig.NewOrganize {
				if err = z.Metadatei.GetEtiketten().EachPtr(
					func(e *kennung.Etikett) (err error) {
						if a.Contains(e) {
							return
						}

						out.Named.AddEtikett(
							z.Metadatei.Bezeichnung.String(),
							*e,
							z.Metadatei.Bezeichnung,
						)

						return
					},
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, c := range a.Children {
		if err = c.addToCompareMap(ot, m, es, out); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
