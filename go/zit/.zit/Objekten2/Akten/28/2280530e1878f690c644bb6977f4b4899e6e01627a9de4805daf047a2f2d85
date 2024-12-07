package organize_text

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
)

type Text struct {
	Options
	*Assignment // TODO make not embedded
}

func New(options Options) (ot *Text, err error) {
	if !options.wasMade {
		panic("options not initialized")
	}

	ot, err = options.Make()

	return
}

func (t *Text) Refine() (err error) {
	if !t.Options.wasMade {
		panic("options not initialized")
	}

	if err = t.Options.refiner().Refine(t.Assignment); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t *Text) ReadFrom(r io.Reader) (n int64, err error) {
	if !t.Options.wasMade {
		panic("options not initialized")
	}

	r1 := &reader{
		options: t.Options,
	}

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

	ocs := t.OptionComments

	for _, oc := range ocs {
		if ocwa, ok := oc.(OptionCommentWithApply); ok {
			if err = ocwa.ApplyToReader(t.Options, r1); err != nil {
				err = errors.Wrapf(err, "OptionComment: %s", oc)
				return
			}
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
	if !ot.Options.wasMade {
		panic("options not initialized")
	}

	lw := format.NewLineWriter()

	omit := ot.HasMetadataContent()

	aw := writer{
		ObjectFactory:        ot.ObjectFactory,
		LineWriter:           lw,
		maxDepth:             ot.MaxDepth(),
		Metadata:             ot.AsMetadata(),
		OmitLeadingEmptyLine: omit,
		options:              ot.Options,
	}

	ocs := ot.OptionComments

	for _, oc := range ocs {
		if ocwa, ok := oc.(OptionCommentWithApply); ok {
			if err = ocwa.ApplyToWriter(ot.Options, &aw); err != nil {
				err = errors.Wrapf(err, "OptionComment: %s", oc)
				return
			}
		}
	}

	if err = aw.write(ot.Assignment); err != nil {
		err = errors.Wrap(err)
		return
	}

	mw := object_metadata.Writer{
		Blob: lw,
	}

	mw.Metadata = ot.Metadata

	if n, err = mw.WriteTo(out); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
