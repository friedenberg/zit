package organize_text

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/line_format"
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

	ot.assignment.addChild(options.AssignmentTreeConstructor.RootAssignment())

	refiner := AssignmentTreeRefiner{
		Enabled:         true,
		UsePrefixJoints: true,
	}

	if err = ot.Refine(refiner); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (t *organizeText) Refine(refiner AssignmentTreeRefiner) (err error) {
	if err = refiner.Refine(t.assignment); err != nil {
		err = errors.Error(err)
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
