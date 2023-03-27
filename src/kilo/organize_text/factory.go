package organize_text

import (
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	zettel_pkg "github.com/friedenberg/zit/src/juliett/zettel"
)

type Factory struct {
	Options
}

func (atc *Factory) Make() (ot *Text, err error) {
	ot = &Text{
		Options:    atc.Options,
		assignment: newAssignment(0),
	}

	ot.assignment.isRoot = true

	ot.Metadatei.EtikettSet = atc.Options.RootEtiketten
	ot.Metadatei.Typ = atc.Options.Typ

	prefixSet := atc.Transacted.ToSetPrefixVerzeichnisse()

	for _, e := range atc.ExtraEtiketten.Elements() {
		ee := newAssignment(ot.Depth() + 1)
		ee.etiketten = kennung.MakeEtikettSet(e)
		ot.assignment.addChild(ee)

		segments := prefixSet.Subset(e)

		var used zettel_pkg.MutableSet

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

func (atc Factory) makeChildren(
	parent *assignment,
	prefixSet zettel_pkg.SetPrefixVerzeichnisse,
	groupingEtiketten kennung.Slice,
) (used zettel_pkg.MutableSet, err error) {
	used = zettel_pkg.MakeMutableSetUnique(0)

	if groupingEtiketten.Len() == 0 {
		prefixSet.ToSet().Each(used.Add)

		err = prefixSet.EachZettel(
			func(e kennung.Etikett, tz zettel_pkg.Transacted) (err error) {
				var z zettel

				if z, err = makeZettel(&tz, atc.Abbr); err != nil {
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
		func(tz *zettel_pkg.Transacted) (err error) {
			var z zettel

			if z, err = makeZettel(tz, atc.Abbr); err != nil {
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
		func(e kennung.Etikett, zs zettel_pkg.MutableSet) (err error) {
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

					if ls, err = e.LeftSubtract(groupingEtiketten[0]); err != nil {
						err = errors.Wrap(err)
						return
					}

					child.etiketten = kennung.MakeEtikettSet(ls)

					nextGroupingEtiketten := kennung.MakeSlice()

					if groupingEtiketten.Len() > 1 {
						nextGroupingEtiketten = groupingEtiketten[1:]
					}

					var usedChild zettel_pkg.MutableSet

					usedChild, err = atc.makeChildren(child, zs.ToSetPrefixVerzeichnisse(), nextGroupingEtiketten)

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

				var usedChild zettel_pkg.MutableSet

				usedChild, err = atc.makeChildren(child, zs.ToSetPrefixVerzeichnisse(), nextGroupingEtiketten)

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
		vi := collections.StringCommaSeparated[kennung.Etikett](parent.children[i].etiketten)
		vj := collections.StringCommaSeparated[kennung.Etikett](parent.children[j].etiketten)
		return vi < vj
	})

	return
}
