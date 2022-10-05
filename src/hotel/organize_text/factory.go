package organize_text

import (
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

type Factory struct {
	Options
}

func (atc *Factory) Make() (ot *Text, err error) {
	ot = &Text{
		Options:    atc.Options,
		assignment: newAssignment(0),
	}

	ot.Metadatei.Set = atc.Options.RootEtiketten
	ot.Metadatei.Typ = atc.Options.Typ

	prefixSet := atc.Transacted.ToSetPrefixTransacted()

	if err = atc.makeChildren(ot.assignment, prefixSet, atc.GroupingEtiketten); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, e := range atc.ExtraEtiketten.Etiketten() {
		child := newAssignment(1)
		child.etiketten = etikett.MakeSet(e)
		ot.assignment.addChild(child)
	}

	if err = ot.Refine(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (atc Factory) makeChildren(
	parent *assignment,
	prefixSet zettel_transacted.SetPrefixTransacted,
	groupingEtiketten etikett.Slice,
) (err error) {
	if groupingEtiketten.Len() == 0 {
		err = prefixSet.EachZettel(
			func(e etikett.Etikett, tz zettel_transacted.Zettel) (err error) {
				var z zettel

				if z, err = makeZettel(tz.Named, atc.Abbr); err != nil {
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
		func(tz zettel_transacted.Zettel) (err error) {
			var z zettel

			if z, err = makeZettel(tz.Named, atc.Abbr); err != nil {
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
		func(e etikett.Etikett, zs zettel_transacted.Set) (err error) {
			if atc.UsePrefixJoints {
				if parent.etiketten.Len() > 1 {
				} else {
					prefixJoint := etikett.MakeSet(groupingEtiketten[0])

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
					child.etiketten = etikett.MakeSet(e.LeftSubtract(groupingEtiketten[0]))

					nextGroupingEtiketten := etikett.NewSlice()

					if groupingEtiketten.Len() > 1 {
						nextGroupingEtiketten = groupingEtiketten[1:]
					}

					err = atc.makeChildren(child, zs.ToSetPrefixTransacted(), nextGroupingEtiketten)

					if err != nil {
						err = errors.Wrap(err)
						return
					}

					intermediate.addChild(child)
				}
			} else {
				child := newAssignment(parent.Depth() + 1)
				child.etiketten = etikett.MakeSet(e)

				nextGroupingEtiketten := etikett.NewSlice()

				if groupingEtiketten.Len() > 1 {
					nextGroupingEtiketten = groupingEtiketten[1:]
				}

				err = atc.makeChildren(child, zs.ToSetPrefixTransacted(), nextGroupingEtiketten)

				if err != nil {
					err = errors.Wrap(err)
					return
				}

				parent.addChild(child)
			}
			return
		},
	)

	sort.Slice(parent.children, func(i, j int) bool {
		return parent.children[i].etiketten.String() < parent.children[j].etiketten.String()
	})

	return
}
