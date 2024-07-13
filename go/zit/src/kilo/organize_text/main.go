package organize_text

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
)

type Text struct {
	Options
	Metadatei
	*Assignment
}

func New(options Options) (ot *Text, err error) {
	if !options.wasMade {
		panic("options not initialized")
	}

	ot, err = options.Make()

	return
}

func (t *Text) Refine() (err error) {
	if err = t.Options.refiner().Refine(t.Assignment); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t *Text) ReadFrom(r io.Reader) (n int64, err error) {
	r1 := &assignmentLineReader{
		options: t.Options,
	}

	r1.stringFormatReader = &t.skuFmt

	mr := object_metadata.Reader{
		Metadatei: &t.Metadatei,
		Akte:      r1,
	}

	var n1 int64
	n1, err = mr.ReadMetadateiFrom(&r)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	ocf := optionCommentFactory{}
	var ocs []Option

	if ocs, err = t.GetOptionComments(ocf); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, oc := range ocs {
		if err = oc.ApplyToReader(t.Options, r1); err != nil {
			err = errors.Wrapf(err, "OptionComment: %s", oc)
			return
		}
	}

	n1, err = mr.ReadBlobFrom(r)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	t.Assignment = r1.root

	return
}

func (ot Text) WriteTo(out io.Writer) (n int64, err error) {
	lw := format.NewLineWriter()

	kopf, schwanz := ot.MaxKopfUndSchwanz(ot.Options)

	l := ot.MaxLen()

	omit := ot.UseMetadateiHeader && ot.HasMetadateiContent()

	aw := assignmentLineWriter{
		LineWriter:           lw,
		maxDepth:             ot.MaxDepth(),
		maxKopf:              kopf,
		maxSchwanz:           schwanz,
		maxLen:               l,
		Metadatei:            ot.AsMetadatei(),
		RightAlignedIndents:  ot.UseRightAlignedIndents,
		OmitLeadingEmptyLine: omit,
	}

	aw.stringFormatWriter = &ot.skuFmt

	ocf := optionCommentFactory{}
	var ocs []Option

	if ocs, err = ot.GetOptionComments(ocf); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, oc := range ocs {
		if err = oc.ApplyToWriter(ot.Options, &aw); err != nil {
			err = errors.Wrapf(err, "OptionComment: %s", oc)
			return
		}
	}

	if aligned, ok := aw.stringFormatWriter.(sku_fmt.KennungAlignedFormat); ok {
		aligned.SetMaxKopfUndSchwanz(kopf, schwanz)
	}

	if err = aw.write(ot.Assignment); err != nil {
		err = errors.Wrap(err)
		return
	}

	mw := object_metadata.Writer{
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
