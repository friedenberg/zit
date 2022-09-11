package organize_text

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/line_format"
)

type Text interface {
	io.ReaderFrom
	io.WriterTo
	ToCompareMap() (out CompareMap, err error)
	Refine(AssignmentTreeRefiner) error
}

type organizeText struct {
	*assignment
}

func New(options Options) (ot *organizeText, err error) {
	ot = &organizeText{
		assignment: newAssignment(),
	}

	ot.assignment.isRoot = true

	var as []*assignment
	as, err = options.AssignmentTreeConstructor.Assignments()

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, a := range as {
		ot.assignment.addChild(a)
	}

	refiner := AssignmentTreeRefiner{
		Enabled:         true,
		UsePrefixJoints: true,
	}

	if err = ot.Refine(refiner); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Print(ot.assignment.etiketten)
	errors.Print(ot.assignment.named)

	return
}

func (t *organizeText) Refine(refiner AssignmentTreeRefiner) (err error) {
	if err = refiner.Refine(t.assignment); err != nil {
		err = errors.Wrap(err)
		return
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

	kopf, scwhanz := ot.assignment.MaxKopfUndSchwanz()

	aw := assignmentLineWriter{
		Writer:              lw,
		maxDepth:            ot.assignment.MaxDepth(),
		maxKopf:             kopf,
		maxScwhanz:          scwhanz,
		experimentalIndents: true,
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
