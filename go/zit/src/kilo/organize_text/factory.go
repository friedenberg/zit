package organize_text

import (
	"sort"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/objekte_collections"
)

type Factory struct {
	Options
}

func (f *Factory) Make() (ot *Text, err error) {
	if f.UseMetadateiHeader {
		ot, err = f.makeWithMetadatei()
	} else {
		ot, err = f.makeWithoutMetadatei()
	}

	return
}

func (atc *Factory) makeWithMetadatei() (ot *Text, err error) {
	ot = &Text{
		Options:    atc.Options,
		Assignment: newAssignment(0),
	}

	ot.IsRoot = true

	ot.EtikettSet = atc.rootEtiketten
	ot.Metadatei.Typ = atc.Typ

	prefixSet := objekte_collections.MakeSetPrefixVerzeichnisse(0)
	atc.Transacted.Each(prefixSet.Add)

	for _, e := range iter.Elements(atc.ExtraEtiketten) {
		ee := newAssignment(ot.GetDepth() + 1)
		ee.Etiketten = kennung.MakeEtikettSet(e)
		ot.addChild(ee)

		segments := prefixSet.Subset(e)

		var used objekte_collections.MutableSetMetadateiWithKennung

		if used, err = atc.makeChildren(
			ee,
			segments.Grouped,
			kennung.MakeEtikettSlice(e),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		prefixSet = prefixSet.Subtract(used)
	}

	if _, err = atc.makeChildren(ot.Assignment, prefixSet, atc.GroupingEtiketten); err != nil {
		err = errors.Wrapf(err, "Assignment: %#v", ot.Assignment)
		return
	}

	if err = ot.Refine(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f Factory) makeWithoutMetadatei() (ot *Text, err error) {
	if !f.wasMade {
		panic("options no initialized")
	}

	ot = &Text{
		Options:    f.Options,
		Assignment: newAssignment(0),
		Metadatei: Metadatei{
			EtikettSet: kennung.MakeEtikettSet(),
		},
	}

	ot.IsRoot = true

	var as []*Assignment
	as, err = f.Options.assignmentTreeConstructor().Assignments()
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, a := range as {
		ot.addChild(a)
	}

	if err = ot.Refine(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (atc Factory) makeChildren(
	parent *Assignment,
	prefixSet objekte_collections.SetPrefixVerzeichnisse,
	groupingEtiketten kennung.EtikettSlice,
) (used objekte_collections.MutableSetMetadateiWithKennung, err error) {
	used = objekte_collections.MakeMutableSetMetadateiWithKennung()

	if groupingEtiketten.Len() == 0 {
		prefixSet.ToSet().Each(used.Add)

		err = prefixSet.EachZettel(
			func(e kennung.Etikett, tz *sku.Transacted) (err error) {
				var z *obj

				if z, err = makeObj(
					atc.PrintOptions,
					tz,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				parent.Named.Add(z)

				return
			},
		)
		if err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	segments := prefixSet.Subset(groupingEtiketten[0])

	err = segments.Ungrouped.Each(
		func(tz *sku.Transacted) (err error) {
			var z *obj

			if z, err = makeObj(
				atc.PrintOptions,
				tz,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			parent.Named.Add(z)

			return
		},
	)
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = segments.Grouped.Each(
		func(e kennung.Etikett, zs objekte_collections.MutableSetMetadateiWithKennung) (err error) {
			if atc.UsePrefixJoints {
				if parent.Etiketten.Len() > 1 {
				} else {
					prefixJoint := kennung.MakeEtikettSet(groupingEtiketten[0])

					var intermediate, lastChild *Assignment

					if len(parent.Children) > 0 {
						lastChild = parent.Children[len(parent.Children)-1]
					}

					if lastChild != nil && (iter.SetEqualsPtr[kennung.Etikett, *kennung.Etikett](lastChild.Etiketten, prefixJoint) || lastChild.Etiketten.Len() == 0) {
						intermediate = lastChild
					} else {
						intermediate = newAssignment(parent.GetDepth() + 1)
						intermediate.Etiketten = prefixJoint
						parent.addChild(intermediate)
					}

					child := newAssignment(intermediate.GetDepth() + 1)

					var ls kennung.Etikett
					b := groupingEtiketten[0]

					if e.Equals(b) {
						return
					}

					if ls, err = kennung.LeftSubtract(e, b); err != nil {
						err = errors.Wrap(err)
						return
					}

					child.Etiketten = kennung.MakeEtikettSet(ls)

					nextGroupingEtiketten := kennung.MakeEtikettSlice()

					if groupingEtiketten.Len() > 1 {
						nextGroupingEtiketten = groupingEtiketten[1:]
					}

					var usedChild objekte_collections.MutableSetMetadateiWithKennung

					psv := objekte_collections.MakeSetPrefixVerzeichnisse(0)
					zs.Each(psv.Add)
					usedChild, err = atc.makeChildren(child, psv, nextGroupingEtiketten)
					if err != nil {
						err = errors.Wrap(err)
						return
					}

					usedChild.Each(used.Add)

					intermediate.addChild(child)
				}
			} else {
				child := newAssignment(parent.GetDepth() + 1)
				child.Etiketten = kennung.MakeEtikettSet(e)

				nextGroupingEtiketten := kennung.MakeEtikettSlice()

				if groupingEtiketten.Len() > 1 {
					nextGroupingEtiketten = groupingEtiketten[1:]
				}

				var usedChild objekte_collections.MutableSetMetadateiWithKennung

				psv := objekte_collections.MakeSetPrefixVerzeichnisse(0)
				zs.Each(psv.Add)
				usedChild, err = atc.makeChildren(child, psv, nextGroupingEtiketten)
				if err != nil {
					err = errors.Wrap(err)
					return
				}

				usedChild.Each(used.Add)

				parent.addChild(child)
			}
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sort.Slice(parent.Children, func(i, j int) bool {
		vi := iter.StringCommaSeparated[kennung.Etikett](
			parent.Children[i].Etiketten,
		)
		vj := iter.StringCommaSeparated[kennung.Etikett](
			parent.Children[j].Etiketten,
		)
		return vi < vj
	})

	return
}
