package object_metadata

import (
	"io"
	"path"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
)

type textParser struct {
	awf interfaces.BlobWriterFactory
	af  script_config.RemoteScript
}

func MakeTextParser(
	awf interfaces.BlobWriterFactory,
	akteFormatter script_config.RemoteScript,
) TextParser {
	if awf == nil {
		panic("nil AkteWriterFactory")
	}

	return textParser{
		awf: awf,
		af:  akteFormatter,
	}
}

func (f textParser) ParseMetadata(
	r io.Reader,
	c TextParserContext,
) (n int64, err error) {
	m := c.GetMetadata()
	Resetter.Reset(m)

	var n1 int64

	defer func() {
		c.SetBlobSha(&m.Blob)
	}()

	var akteFD fd.FD

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
						return f.readTyp(m, v, &akteFD)
					},
				},
			),
		),
	)

	var akteWriter sha.WriteCloser

	if akteWriter, err = f.awf.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if akteWriter == nil {
		err = errors.Errorf("akte writer is nil")
		return
	}

	defer errors.DeferredCloser(&err, akteWriter)

	mr := Reader{
		// RequireMetadatei: true,
		Metadata: lr,
		Blob:      akteWriter,
	}

	// if cmg, ok := c.(checkout_mode.Getter); ok {
	// 	var checkoutMode checkout_mode.Mode

	// 	if checkoutMode, err = cmg.GetCheckoutMode(); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// 	mr.RequireMetadatei = checkoutMode.IncludesObjekte()
	// }

	if n, err = mr.ReadFrom(r); err != nil {
		n += n1
		err = errors.Wrap(err)
		return
	}

	n += n1

	inlineAkteSha := sha.Make(akteWriter.GetShaLike())

	if !m.Blob.IsNull() && !akteFD.GetShaLike().IsNull() {
		err = errors.Wrap(
			MakeErrHasInlineBlobAndFilePath(
				&akteFD,
				inlineAkteSha,
			),
		)

		return
	} else if !akteFD.GetShaLike().IsNull() {
		if afs, ok := c.(fd.AkteFDSetter); ok {
			afs.SetAkteFD(&akteFD)
		}

		m.Blob.SetShaLike(akteFD.GetShaLike())
	}

	switch {
	case m.Blob.IsNull() && !inlineAkteSha.IsNull():
		m.Blob.SetShaLike(inlineAkteSha)

	case !m.Blob.IsNull() && inlineAkteSha.IsNull():
		// noop

	case !m.Blob.IsNull() && !inlineAkteSha.IsNull() &&
		!m.Blob.Equals(inlineAkteSha):
		err = errors.Wrap(
			MakeErrHasInlineBlobAndMetadateiSha(
				inlineAkteSha,
				&m.Blob,
			),
		)

		return
	}

	return
}

func (f textParser) readTyp(
	m *Metadata,
	desc string,
	akteFD *fd.FD,
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

		if err = akteFD.SetWithAkteWriterFactory(
			desc,
			f.awf,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	//! <sha>.<typ ext>
	case tail != "":
		if err = f.setAkteSha(m, head); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = m.Type.Set(tail); err != nil {
			err = errors.Wrap(err)
			return
		}

	//! <sha>
	case tail == "":
		if err = f.setAkteSha(m, head); err == nil {
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

func (f textParser) setAkteSha(
	m *Metadata,
	maybeSha string,
) (err error) {
	if err = m.Blob.Set(maybeSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
