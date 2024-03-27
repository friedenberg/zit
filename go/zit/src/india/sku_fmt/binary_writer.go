package sku_fmt

import (
	"bytes"
	"encoding"
	"fmt"
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/catgut"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/delta/ohio"
	"code.linenisgreat.com/zit/src/delta/schlussel"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

type BinaryWriter struct {
	bytes.Buffer
	BinaryField
	kennung.Sigil
}

func (bf *BinaryWriter) UpdateSigil(
	wa io.WriterAt,
	s kennung.Sigil,
	offset int64,
) (err error) {
	s.Add(bf.Sigil)
	// 2 uint8 + offset + 2 uint8 + Schlussel
	offset = int64(2) + offset + int64(3)

	var n int

	if n, err = wa.WriteAt([]byte{s.Byte()}, offset); err != nil {
		err = errors.Wrap(err)
		return
	}

	if n != 1 {
		err = catgut.MakeErrLength(1, int64(n), nil)
		return
	}

	return
}

func (bf *BinaryWriter) WriteFormat(
	w io.Writer,
	sk SkuWithSigil,
) (n int64, err error) {
	bf.Buffer.Reset()

	for _, f := range binaryFieldOrder {
		bf.BinaryField.Reset()
		bf.Schlussel = f

		if _, err = bf.writeFieldKey(sk); err != nil {
			err = errors.Wrapf(err, "Sku: %s", sk)
			return
		}
	}

	bf.BinaryField.Reset()

	bf.SetContentLength(bf.Len())

	var n1 int
	n1, err = ohio.WriteAllOrDieTrying(w, bf.ContentLength[:])
	n += int64(n1)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	var n2 int64
	n2, err = bf.Buffer.WriteTo(w)
	n += n2

	return
}

func (bf *BinaryWriter) writeFieldKey(
	sk SkuWithSigil,
) (n int64, err error) {
	switch bf.Schlussel {
	case schlussel.Sigil:
		s := sk.Sigil
		s.Add(bf.Sigil)

		if sk.Metadatei.Verzeichnisse.Archiviert.Bool() {
			s.Add(kennung.SigilHidden)
		}

		if n, err = bf.writeFieldByteReader(s); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.Akte:
		if n, err = bf.writeSha(&sk.Metadatei.Akte, true); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.Bezeichnung:
		if sk.Metadatei.Bezeichnung.IsEmpty() {
			return
		}

		if n, err = bf.writeFieldBinaryMarshaler(&sk.Metadatei.Bezeichnung); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.Etikett:
		es := sk.GetEtiketten()

		for _, e := range iter.SortedValues[kennung.Etikett](es) {
			if e.IsVirtual() {
				continue
			}

			var n1 int64
			n1, err = bf.writeFieldBinaryMarshaler(&e)
			n += n1

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	case schlussel.Kennung:
		if n, err = bf.writeFieldWriterTo(&sk.Kennung); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.Tai:
		if n, err = bf.writeFieldWriterTo(&sk.Metadatei.Tai); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.Typ:
		if sk.Metadatei.Typ.IsEmpty() {
			return
		}

		if n, err = bf.writeFieldBinaryMarshaler(&sk.Metadatei.Typ); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.MutterMetadateiMutterKennung:
		if n, err = bf.writeSha(sk.Metadatei.Mutter(), true); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.ShaMetadateiSansTai:
		if n, err = bf.writeSha(&sk.Metadatei.SelbstMetadateiSansTai, false); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.ShaMetadateiMutterKennung:
		if n, err = bf.writeSha(sk.Metadatei.Sha(), false); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.ShaMetadatei:
		if err = sha.MakeErrIsNull(&sk.Metadatei.SelbstMetadatei); err != nil {
			return
		}

		if n, err = bf.writeFieldWriterTo(&sk.Metadatei.SelbstMetadatei); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.VerzeichnisseEtikettImplicit:
		es := sk.Metadatei.Verzeichnisse.GetImplicitEtiketten()

		for _, e := range iter.SortedValues[kennung.Etikett](es) {
			var n1 int64
			n1, err = bf.writeFieldBinaryMarshaler(&e)
			n += n1

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	case schlussel.VerzeichnisseEtikettExpanded:
		es := sk.Metadatei.Verzeichnisse.GetExpandedEtiketten()

		for _, e := range iter.SortedValues[kennung.Etikett](es) {
			var n1 int64
			n1, err = bf.writeFieldBinaryMarshaler(&e)
			n += n1

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	case schlussel.VerzeichnisseEtiketten:
		es := sk.Metadatei.Verzeichnisse.Etiketten

		for _, e := range es {
			var n1 int64
			n1, err = bf.writeFieldWriterTo(e)
			n += n1

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	default:
		panic(fmt.Sprintf("unsupported key: %s", bf.Schlussel))
	}

	return
}

func (bf *BinaryWriter) writeSha(
	sh *sha.Sha,
	allowNull bool,
) (n int64, err error) {
	if sh.IsNull() {
		if !allowNull {
			err = errors.Wrap(sha.ErrIsNull)
		}

		return
	}

	if n, err = bf.writeFieldWriterTo(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (bf *BinaryWriter) writeFieldWriterTo(
	wt io.WriterTo,
) (n int64, err error) {
	_, err = wt.WriteTo(&bf.Content)

	if err != nil {
		return
	}

	if n, err = bf.BinaryField.WriteTo(&bf.Buffer); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (bf *BinaryWriter) writeFieldBinaryMarshaler(
	bm encoding.BinaryMarshaler,
) (n int64, err error) {
	var b []byte

	b, err = bm.MarshalBinary()

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = ohio.WriteAllOrDieTrying(&bf.Content, b); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err != nil {
		err = errors.WrapExceptAsNil(err, io.EOF)
		return
	}

	if n, err = bf.BinaryField.WriteTo(&bf.Buffer); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (bf *BinaryWriter) writeFieldByteReader(
	br io.ByteReader,
) (n int64, err error) {
	var b byte

	b, err = br.ReadByte()

	if err != nil {
		return
	}

	err = bf.Content.WriteByte(b)

	if err != nil {
		return
	}

	if n, err = bf.BinaryField.WriteTo(&bf.Buffer); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
