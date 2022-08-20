package organize_text

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/line_format"
)

type Text interface {
	io.ReaderFrom
	io.WriterTo
	ToCompareMap() (out CompareMap)
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
		UsePrefixJoints: true,
	}

	if err = refiner.Refine(ot.assignment); err != nil {
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
