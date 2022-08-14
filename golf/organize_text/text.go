package organize_text

import (
	"io"
	"sort"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/charlie/line_format"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

type Text interface {
	io.ReaderFrom
	io.WriterTo
	ToCompareMap() (out CompareMap)
}

type organizeText struct {
	*assignment
}

func New(options Options, named stored_zettel.SetNamed) (ot *organizeText, err error) {
	ot = NewEmpty()

	root := newAssignment(1)
	root.etiketten = options.RootEtiketten
	// root.named = makeZettelZetFromSetNamed(named)

	prefixSet := named.ToSetPrefixNamed()
	makeChildren(root, *prefixSet, options.GroupingEtiketten)

	for _, e := range options.ExtraEtiketten {
		child := newAssignment(root.depth + 1)
		child.etiketten = etikett.MakeSet(e)
		root.addChild(child)
	}

	ot.assignment.addChild(root)

	return
}

func makeChildren(
	parent *assignment,
	prefixSet stored_zettel.SetPrefixNamed,
	remainingEtiketten etikett.Slice) (assigned *stored_zettel.SetNamed) {

	assigned = stored_zettel.NewSetNamed()

	// logz.Print("making children")
	if remainingEtiketten.Len() == 0 {
		for _, zs := range prefixSet {
			// assigned.Merge(zs)
			for _, z := range zs {
				// logz.Printf("%s adding named %s", parent.etiketten, z.Hinweis)
				parent.named.Add(makeZettel(z))
			}
		}

		return
	}

	segments := prefixSet.Subset(remainingEtiketten[0])
	// logz.Printf("head: %s ungrouped: %s", remainingEtiketten[0], segments.Ungrouped.HinweisStrings())
	// logz.Printf("head: %s grouped: %s", remainingEtiketten[0], segments.Grouped.ToSetNamed().HinweisStrings())

	for _, z := range *segments.Ungrouped {
		parent.named.Add(makeZettel(z))
	}

	for e, zs := range *segments.Grouped {
		// assigned.Merge(zs)
		// logz.Print("iterating through grouped: ", e)
		child := newAssignment(parent.depth + 1)
		child.etiketten = etikett.MakeSet(e)
		// child.named = makeZettelZetFromSetNamed(zs)

		nextEtiketten := etikett.NewSlice()

		if remainingEtiketten.Len() > 1 {
			nextEtiketten = remainingEtiketten[1:]
		}

		_ = makeChildren(child, *zs.ToSetPrefixNamed(), nextEtiketten)
		// childAssigned.Merge(c)
		// assigned.Merge(*c)

		parent.addChild(child)

		sort.Slice(parent.children, func(i, j int) bool {
			return parent.children[i].etiketten.String() < parent.children[j].etiketten.String()
		})
	}

	return
}

func NewEmpty() (ot *organizeText) {
	ot = &organizeText{
		assignment: newAssignment(0),
	}

	return
}

func (t *organizeText) ReadFrom(r io.Reader) (n int64, err error) {
	r1 := assignmentLineReader{}

	n, err = r1.ReadFrom(r)

	t.assignment = r1.root

	return
}

func (ot organizeText) WriteTo(out io.Writer) (n int64, err error) {
	lw := line_format.NewWriter()

	aw := assignmentLineWriter{Writer: lw}

	if err = aw.write(ot.assignment); err != nil {
		err = errors.Error(err)
		return
	}

	if n, err = lw.WriteTo(out); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
