package organize_text

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/format"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/india/sku_fmt"
)

type Text struct {
	Options
	Metadatei
	*assignment
}

func New(options Options) (ot *Text, err error) {
	if !options.wasMade {
		panic("options not initialized")
	}

	ot, err = options.Factory().Make()

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
	r1 := &assignmentLineReader{
		options: t.Options,
		ex:      t.Expanders,
	}

	mr := metadatei.Reader{
		Metadatei: &t.Metadatei,
		Akte:      r1,
	}

	n, err = mr.ReadFrom(r)

	t.assignment = r1.root

	return
}

func (ot Text) WriteTo(out io.Writer) (n int64, err error) {
	lw := format.NewLineWriter()

	kopf, schwanz := ot.MaxKopfUndSchwanz()

	l := ot.MaxLen()

	omit := ot.UseMetadateiHeader && ot.HasMetadateiContent()

	aw := assignmentLineWriter{
		LineWriter:           lw,
		maxDepth:             ot.MaxDepth(),
		maxKopf:              kopf,
		maxSchwanz:           schwanz,
		maxLen:               l,
		RightAlignedIndents:  ot.UseRightAlignedIndents,
		OmitLeadingEmptyLine: omit,
		stringFormatWriter:   sku_fmt.MakeOrganizeFormat(kopf, schwanz),
	}

	if err = aw.write(ot.assignment); err != nil {
		err = errors.Wrap(err)
		return
	}

	mw := metadatei.Writer{
		Akte: lw,
	}

	if ot.UseMetadateiHeader {
		ot.Matchers = ot.commentMatchers
		mw.Metadatei = ot.Metadatei
	}

	if n, err = mw.WriteTo(out); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
