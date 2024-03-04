package organize_text2

import (
	"fmt"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/lima/changes2"
)

func (ot *Text) CompareMap(
	hinweis_expander func(string) (*kennung.Hinweis, error),
) (out changes2.CompareMap, err error) {
	preExpansion := changes2.CompareMap{
		Named:   make(changes2.SetKeyToMetadatei),
		Unnamed: make(changes2.SetKeyToMetadatei),
	}

	if err = ot.addToCompareMap(
		ot,
		ot.Metadatei,
		kennung.MakeEtikettSet(),
		&preExpansion,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	out = changes2.CompareMap{
		Named:   make(changes2.SetKeyToMetadatei),
		Unnamed: preExpansion.Unnamed,
	}

	for h, v := range preExpansion.Named {
		var h1 schnittstellen.Stringer

		if h1, err = hinweis_expander(h); err == nil {
			h = h1.String()
		}

		err = nil

		out.Named[h] = v
	}

	return
}

func (a *Assignment) addToCompareMap(
	ot *Text,
	m Metadatei,
	es kennung.EtikettSet,
	out *changes2.CompareMap,
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

func (ot *Text) GetSkus() sku.TransactedSet {
	out := sku.MakeTransactedMutableSet()

	return out
}
