package organize_text

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/line_format"
	"github.com/friedenberg/zit/src/bravo/metadatei_io"
	"github.com/friedenberg/zit/src/delta/etikett"
)

type Text struct {
	Options
	Metadatei
	*assignment
}

func New(options Options) (ot *Text, err error) {
	if !options.wasMade {
		panic("options no initialized")
	}

	if options.UseMetadateiHeader {
		ot, err = options.Factory().Make()
	} else {
		ot, err = newWithoutMetadatei(options)
	}

	return
}

func newWithoutMetadatei(options Options) (ot *Text, err error) {
	if !options.wasMade {
		panic("options no initialized")
	}

	ot = &Text{
		Options:    options,
		assignment: newAssignment(0),
		Metadatei: Metadatei{
			Set: etikett.MakeSet(),
		},
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
	r1 := &assignmentLineReader{}

	mr := metadatei_io.Reader{
		Metadatei: &t.Metadatei,
		Akte:      r1,
	}

	n, err = mr.ReadFrom(r)

	t.assignment = r1.root

	return
}

func (ot Text) WriteTo(out io.Writer) (n int64, err error) {
	lw := line_format.NewWriter()

	kopf, scwhanz := ot.assignment.MaxKopfUndSchwanz()

	aw := assignmentLineWriter{
		Writer:               lw,
		maxDepth:             ot.assignment.MaxDepth(),
		maxKopf:              kopf,
		maxScwhanz:           scwhanz,
		RightAlignedIndents:  ot.UseRightAlignedIndents,
		OmitLeadingEmptyLine: ot.Options.UseMetadateiHeader && ot.Metadatei.HasMetadateiContent(),
	}

	if err = aw.write(ot.assignment); err != nil {
		err = errors.Wrap(err)
		return
	}

	mw := metadatei_io.Writer{
		Akte: lw,
	}

	if ot.Options.UseMetadateiHeader {
		mw.Metadatei = ot.Metadatei
	}

	if n, err = mw.WriteTo(out); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}