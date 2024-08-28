package organize_text

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
)

type Text struct {
	Options
	Metadata
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

	r1.stringFormatReader = t.stringFormatReader

	mr := object_metadata.Reader{
		Metadata: &t.Metadata,
		Blob:     r1,
	}

	var n1 int64
	n1, err = mr.ReadMetadataFrom(&r)
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

	kopf, schwanz := ot.MaxHeadAndTail(ot.Options)

	l := ot.MaxLen()

	omit := ot.UseMetadataHeader && ot.HasMetadataContent()

	aw := assignmentLineWriter{
		SkuPool:              ot.SkuPool,
		LineWriter:           lw,
		maxDepth:             ot.MaxDepth(),
		maxHead:              kopf,
		maxTail:              schwanz,
		maxLen:               l,
		Metadata:             ot.AsMetadatei(),
		RightAlignedIndents:  ot.UseRightAlignedIndents,
		OmitLeadingEmptyLine: omit,
	}

	aw.stringFormatWriter = ot.stringFormatWriter

	ocf := optionCommentFactory{}
	var ocs []Option

	if ocs, err = ot.GetOptionComments(ocf); err != nil {
		err = errors.Wrap(err)
		return
	}

	if ot.Config.DryRun {
		ocs = append(ocs, optionCommentDryRun(values.MakeBool(true)))
	}

	for _, oc := range ocs {
		if err = oc.ApplyToWriter(ot.Options, &aw); err != nil {
			err = errors.Wrapf(err, "OptionComment: %s", oc)
			return
		}
	}

	if aligned, ok := aw.stringFormatWriter.(sku_fmt.ObjectIdAlignedFormat); ok {
		aligned.SetMaxKopfUndSchwanz(kopf, schwanz)
	}

	if err = aw.write(ot.Assignment); err != nil {
		err = errors.Wrap(err)
		return
	}

	mw := object_metadata.Writer{
		Blob: lw,
	}

	if ot.UseMetadataHeader {
		ot.Matchers = ot.commentMatchers
		mw.Metadata = ot.Metadata
	}

	if n, err = mw.WriteTo(out); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
