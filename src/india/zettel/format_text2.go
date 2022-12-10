package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/metadatei_io"
	"github.com/friedenberg/zit/src/echo/sha"
)

type Text2 struct {
	af gattung.AkteIOFactory

	DoNotWriteEmptyBezeichnung bool
	TypError                   error
}

func (f Text2) ReadFormat(
	r io.Reader,
	z *Objekte,
) (n int64, err error) {
	formatMetadatei := &FormatMetadateiText2{
		af: f.af,
	}

	var akteWriter sha.WriteCloser

	if f.af == nil {
		err = errors.Errorf("akte writer factory is nil")
		return
	}

	if akteWriter, err = f.af.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if akteWriter == nil {
		err = errors.Errorf("akte writer is nil")
		return
	}

	mr := metadatei_io.Reader{
		Metadatei: format.MakeReaderFrom[Objekte](
			formatMetadatei.ReadFormat,
			z,
		),
		Akte: akteWriter,
	}

	if n, err = mr.ReadFrom(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = akteWriter.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	inlineAkteSha := akteWriter.Sha()

	switch {
	case formatMetadatei.akteSha.IsNull() && !inlineAkteSha.IsNull():
		z.Akte = inlineAkteSha

	case !formatMetadatei.akteSha.IsNull() && inlineAkteSha.IsNull():
		z.Akte = formatMetadatei.akteSha

	case !formatMetadatei.akteSha.IsNull() && !inlineAkteSha.IsNull():
		err = ErrHasInlineAkteAndFilePath{
			External: externalFile{
				Sha:  formatMetadatei.akteSha,
				Path: formatMetadatei.aktePath,
			},
			InlineSha: inlineAkteSha,
			Objekte:   *z,
		}

		return
	}

	return
}

func (f Text2) WriteFormat(
	w io.Writer,
	z *Objekte,
) (n int64, err error) {
	// formatMetadatei := &FormatMetadateiText2{
	// 	af: f.af,
	// }

	// var akteWriter sha.ReadCloser

	// mr := metadatei_io.Writer{
	// 	Metadatei: format.MakeWriterTo2(
	// 		formatMetadatei.WriteFormat,
	// 		z,
	// 	),
	// 	Akte: akteWriter,
	// }

	// switch {
	// case f.IncludeAkte && c.ExternalAktePath == "":
	// 	return f.writeToInlineAkte(c)

	// case c.IncludeAkte:
	// 	return f.writeToExternalAkte(c)

	// default:
	// 	return f.writeToOmitAkte(c)
	// }

	return
}
