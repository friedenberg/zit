package organize_text

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/line_format"
)

type Text struct {
	Options
	*assignment
}

func New(options Options) (ot *Text, err error) {
	ot = &Text{
		Options:    options,
		assignment: newAssignment(),
	}

	ot.assignment.isRoot = true

	var as []*assignment
	as, err = options.assignmentTreeConstructor().Assignments()

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

func (t *Text) Refine() (err error) {
	if err = t.Options.refiner().Refine(t.assignment); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t *Text) ReadFrom(r io.Reader) (n int64, err error) {
	r1 := assignmentLineReader{}

	n, err = r1.ReadFrom(r)

	t.assignment = r1.root

	return
}

func (ot Text) WriteTo(out io.Writer) (n int64, err error) {
	lw := line_format.NewWriter()

	kopf, scwhanz := ot.assignment.MaxKopfUndSchwanz()

	aw := assignmentLineWriter{
		Writer:              lw,
		maxDepth:            ot.assignment.MaxDepth(),
		maxKopf:             kopf,
		maxScwhanz:          scwhanz,
		RightAlignedIndents: ot.UseRightAlignedIndents,
	}

	if err = aw.write(ot.assignment); err != nil {
		err = errors.Wrap(err)
		return
	}

	if n, err = lw.WriteTo(out); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
