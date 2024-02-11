package organize_text

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/echo/format"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
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

	ot, err = options.Factory().Make()

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

	if t.Konfig.NewOrganize {
		r1.stringFormatReader = &t.organizeNew
	} else {
		r1.stringFormatReader = &t.organize
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

	mr := metadatei.Reader{
		Metadatei: &t.Metadatei,
		Akte:      r1,
	}

	n, err = mr.ReadFrom(r)

	t.Assignment = r1.root

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
		Metadatei:            ot.AsMetadatei(),
		RightAlignedIndents:  ot.UseRightAlignedIndents,
		OmitLeadingEmptyLine: omit,
	}

	if ot.Konfig.NewOrganize {
		aw.stringFormatWriter = &ot.organizeNew
	} else {
		aw.stringFormatWriter = &ot.organize
	}

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
