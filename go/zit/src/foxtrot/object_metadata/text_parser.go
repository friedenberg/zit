package object_metadata

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type textParser struct {
	awf interfaces.BlobWriterFactory
	af  script_config.RemoteScript
}

func MakeTextParser(
	awf interfaces.BlobWriterFactory,
	blobFormatter script_config.RemoteScript,
) TextParser {
	if awf == nil {
		panic("nil BlobWriterFactory")
	}

	return textParser{
		awf: awf,
		af:  blobFormatter,
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

	mp := &textParser2{
		BlobWriterFactory: f.awf,
		TextParserContext: c,
	}

	var blobWriter sha.WriteCloser

	if blobWriter, err = f.awf.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if blobWriter == nil {
		err = errors.Errorf("blob writer is nil")
		return
	}

	defer errors.DeferredCloser(&err, blobWriter)

	mr := Reader{
		// RequireMetadatei: true,
		Metadata: mp,
		Blob:     blobWriter,
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

	inlineBlobSha := sha.Make(blobWriter.GetShaLike())

	if !m.Blob.IsNull() && !mp.Blob.GetShaLike().IsNull() {
		err = errors.Wrap(
			MakeErrHasInlineBlobAndFilePath(
				&mp.Blob,
				inlineBlobSha,
			),
		)

		return
	} else if !mp.Blob.GetShaLike().IsNull() {
		if afs, ok := c.(BlobFDSetter); ok {
			afs.SetBlobFD(&mp.Blob)
		}

		m.Blob.SetShaLike(mp.Blob.GetShaLike())
	}

	switch {
	case m.Blob.IsNull() && !inlineBlobSha.IsNull():
		m.Blob.SetShaLike(inlineBlobSha)

	case !m.Blob.IsNull() && inlineBlobSha.IsNull():
		// noop

	case !m.Blob.IsNull() && !inlineBlobSha.IsNull() &&
		!m.Blob.Equals(inlineBlobSha):
		err = errors.Wrap(
			MakeErrHasInlineBlobAndMetadateiSha(
				inlineBlobSha,
				&m.Blob,
			),
		)

		return
	}

	return
}
