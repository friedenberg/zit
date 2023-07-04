package organize_text

import (
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/india/objekte_collections"
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
		assignment: newAssignment(0),
	}

	ot.assignment.isRoot = true

	ot.Metadatei.EtikettSet = atc.Options.RootEtiketten
	ot.Metadatei.Typ = atc.Options.Typ

	prefixSet := objekte_collections.MakeSetPrefixVerzeichnisse(0)
	atc.Transacted.Each(prefixSet.Add)

	for _, e := range atc.ExtraEtiketten.Elements() {
		ee := newAssignment(ot.Depth() + 1)
		ee.etiketten = kennung.MakeEtikettSet(e)
		ot.assignment.addChild(ee)

		segments := prefixSet.Subset(e)

		var used objekte_collections.MutableSetMetadateiWithKennung

		if used, err = atc.makeChildren(ee, segments.Grouped, kennung.MakeSlice(e)); err != nil {
			err = errors.Wrap(err)
			return
		}

		prefixSet = prefixSet.Subtract(used)
	}

	if _, err = atc.makeChildren(ot.assignment, prefixSet, atc.GroupingEtiketten); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = ot.Refine(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f Factory) makeWithoutMetadatei() (ot *Text, err error) {
	if !f.Options.wasMade {
		panic("options no initialized")
	}

	ot = &Text{
		Options:    f.Options,
		assignment: newAssignment(0),
		Metadatei: Metadatei{
			EtikettSet: kennung.MakeEtikettSet(),
		},
	}

	ot.assignment.isRoot = true

	var as []*assignment
	as, err = f.Options.assignmentTreeConstructor().Assignments()

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, a := range as {
		ot.assignment.addChild(a)
	}

	if err = ot.Refine(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (atc Factory) makeChildren(
	parent *assignment,
	prefixSet objekte_collections.SetPrefixVerzeichnisse,
	groupingEtiketten kennung.Slice,
) (used objekte_collections.MutableSetMetadateiWithKennung, err error) {
	used = objekte_collections.MakeMutableSetMetadateiWithKennung()

	if groupingEtiketten.Len() == 0 {
		prefixSet.ToSet().Each(used.Add)

		err = prefixSet.EachZettel(
			func(e kennung.Etikett, tz sku.SkuLike) (err error) {
				var z obj

				if z, err = makeObj(tz, atc.Expanders); err != nil {
					err = errors.Wrap(err)
					return
				}

				parent.named.Add(z)

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
		func(tz sku.SkuLike) (err error) {
			var z obj

			if z, err = makeObj(tz, atc.Expanders); err != nil {
				err = errors.Wrap(err)
				return
			}

			parent.named.Add(z)
			return
		},
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	segments.Grouped.Each(
		func(e kennung.Etikett, zs objekte_collections.MutableSetMetadateiWithKennung) (err error) {
			if atc.UsePrefixJoints {
				if parent.etiketten.Len() > 1 {
				} else {
					prefixJoint := kennung.MakeEtikettSet(groupingEtiketten[0])

					var intermediate, lastChild *assignment

					if len(parent.children) > 0 {
						lastChild = parent.children[len(parent.children)-1]
					}

					if lastChild != nil && (lastChild.etiketten.Equals(prefixJoint) || lastChild.etiketten.Len() == 0) {
						intermediate = lastChild
					} else {
						intermediate = newAssignment(parent.Depth() + 1)
						intermediate.etiketten = prefixJoint
						parent.addChild(intermediate)
					}

					child := newAssignment(intermediate.Depth() + 1)

					var ls kennung.Etikett

					if ls, err = kennung.LeftSubtract(
						e,
						groupingEtiketten[0],
					); err != nil {
						err = errors.Wrap(err)
						return
					}

					child.etiketten = kennung.MakeEtikettSet(ls)

					nextGroupingEtiketten := kennung.MakeSlice()

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
				child := newAssignment(parent.Depth() + 1)
				child.etiketten = kennung.MakeEtikettSet(e)

				nextGroupingEtiketten := kennung.MakeSlice()

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
	)

	sort.Slice(parent.children, func(i, j int) bool {
		vi := collections.StringCommaSeparated[kennung.Etikett](
			parent.children[i].etiketten,
		)
		vj := collections.StringCommaSeparated[kennung.Etikett](
			parent.children[j].etiketten,
		)
		return vi < vj
	})

	return
}
