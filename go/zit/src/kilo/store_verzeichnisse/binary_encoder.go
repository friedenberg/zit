package store_verzeichnisse

import (
	"bytes"
	"encoding"
	"fmt"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/keys"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type binaryEncoder struct {
	bytes.Buffer
	binaryField
	ids.Sigil
}

func (bf *binaryEncoder) updateSigil(
	wa io.WriterAt,
	s ids.Sigil,
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

func (bf *binaryEncoder) writeFormat(
	w io.Writer,
	sk skuWithSigil,
) (n int64, err error) {
	bf.Buffer.Reset()

	for _, f := range binaryFieldOrder {
		bf.binaryField.Reset()
		bf.Key = f

		if _, err = bf.writeFieldKey(sk); err != nil {
			err = errors.Wrapf(err, "Sku: %s", sk)
			return
		}
	}

	bf.binaryField.Reset()

	defer func() {
		r := recover()

		if r == nil {
			return
		}

		ui.Debug().Print(sk, bf.Len(), &sk.Metadata.Cache.TagPaths)
		panic(r)
	}()
	// TODO
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

func (bf *binaryEncoder) writeFieldKey(
	sk skuWithSigil,
) (n int64, err error) {
	switch bf.Key {
	case keys.Sigil:
		s := sk.Sigil
		s.Add(bf.Sigil)

		if sk.Metadata.Cache.Dormant.Bool() {
			s.Add(ids.SigilHidden)
		}

		if n, err = bf.writeFieldByteReader(s); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.Blob:
		if n, err = bf.writeSha(&sk.Metadata.Blob, true); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.Description:
		if sk.Metadata.Description.IsEmpty() {
			return
		}

		if n, err = bf.writeFieldBinaryMarshaler(&sk.Metadata.Description); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.Tag:
		es := sk.GetTags()

		for _, e := range iter.SortedValues(es) {
			if e.IsVirtual() {
				continue
			}

			if e.String() == "" {
				err = errors.Errorf("empty etikett in %q", es)
				return
			}

			var n1 int64
			n1, err = bf.writeFieldBinaryMarshaler(&e)
			n += n1

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	case keys.ObjectId:
		if n, err = bf.writeFieldWriterTo(&sk.ObjectId); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.Tai:
		if n, err = bf.writeFieldWriterTo(&sk.Metadata.Tai); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.Type:
		if sk.Metadata.Type.IsEmpty() {
			return
		}

		if n, err = bf.writeFieldBinaryMarshaler(&sk.Metadata.Type); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.MutterMetadateiMutterKennung:
		if n, err = bf.writeSha(sk.Metadata.Mutter(), true); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.ShaMetadateiSansTai:
		if n, err = bf.writeSha(&sk.Metadata.SelfMetadataWithoutTai, true); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.ShaMetadateiMutterKennung:
		if n, err = bf.writeSha(sk.Metadata.Sha(), false); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.ShaMetadatei:
		if err = sha.MakeErrIsNull(&sk.Metadata.SelfMetadata); err != nil {
			return
		}

		if n, err = bf.writeFieldWriterTo(&sk.Metadata.SelfMetadata); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.VerzeichnisseEtikettImplicit:
		es := sk.Metadata.Cache.GetImplicitTags()

		for _, e := range iter.SortedValues[ids.Tag](es) {
			var n1 int64
			n1, err = bf.writeFieldBinaryMarshaler(&e)
			n += n1

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	case keys.VerzeichnisseEtikettExpanded:
		es := sk.Metadata.Cache.GetExpandedTags()

		for _, e := range iter.SortedValues[ids.Tag](es) {
			var n1 int64
			n1, err = bf.writeFieldBinaryMarshaler(&e)
			n += n1

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	case keys.VerzeichnisseEtiketten:
		es := sk.Metadata.Cache.TagPaths

		for _, e := range es.Paths {
			var n1 int64
			n1, err = bf.writeFieldWriterTo(e)
			n += n1

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	default:
		panic(fmt.Sprintf("unsupported key: %s", bf.Key))
	}

	return
}

func (bf *binaryEncoder) writeSha(
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

func (bf *binaryEncoder) writeFieldWriterTo(
	wt io.WriterTo,
) (n int64, err error) {
	_, err = wt.WriteTo(&bf.Content)
	if err != nil {
		return
	}

	if n, err = bf.binaryField.WriteTo(&bf.Buffer); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (bf *binaryEncoder) writeFieldBinaryMarshaler(
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

	if n, err = bf.binaryField.WriteTo(&bf.Buffer); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (bf *binaryEncoder) writeFieldByteReader(
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

	if n, err = bf.binaryField.WriteTo(&bf.Buffer); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
