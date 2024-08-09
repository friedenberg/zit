package object_metadata

import (
	"io"
	"path"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
)

type textParser2 struct {
	interfaces.BlobWriterFactory
	TextParserContext
	Blob fd.FD
}

func (f *textParser2) ReadFrom(r io.Reader) (n int64, err error) {
	m := f.GetMetadata()
	Resetter.Reset(m)

	lr := format.MakeLineReaderConsumeEmpty(
		ohio.MakeLineReaderIterate(
			ohio.MakeLineReaderKeyValues(
				map[string]interfaces.FuncSetString{
					"#": m.Description.Set,
					"%": func(v string) (err error) {
						m.Comments = append(m.Comments, v)
						return
					},
					"-": m.AddTagString,
					"!": func(v string) (err error) {
						return f.readTyp(m, v)
					},
				},
			),
		),
	)

	if n, err = lr.ReadFrom(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *textParser2) readTyp(
	m *Metadata,
	desc string,
) (err error) {
	if desc == "" {
		return
	}

	tail := path.Ext(desc)
	head := desc[:len(desc)-len(tail)]

	//! <path>.<typ ext>
	switch {
	case files.Exists(desc):
		if err = m.Type.Set(tail); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = f.Blob.SetWithBlobWriterFactory(desc, f); err != nil {
			err = errors.Wrap(err)
			return
		}

	//! <sha>.<typ ext>
	case tail != "":
		if err = f.setBlobSha(m, head); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = m.Type.Set(tail); err != nil {
			err = errors.Wrap(err)
			return
		}

	//! <sha>
	case tail == "":
		if err = f.setBlobSha(m, head); err == nil {
			return
		}

		err = nil

		fallthrough

	//! <typ ext>
	default:
		if err = m.Type.Set(head); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (f *textParser2) setBlobSha(
	m *Metadata,
	maybeSha string,
) (err error) {
	if err = m.Blob.Set(maybeSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
