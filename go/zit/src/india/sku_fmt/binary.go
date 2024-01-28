package sku_fmt

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/ohio"
	"github.com/friedenberg/zit/src/delta/schlussel"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/etiketten_path"
	"github.com/friedenberg/zit/src/golf/ennui"
	"github.com/friedenberg/zit/src/hotel/sku"
)

var binaryFieldOrder = []schlussel.Schlussel{
	schlussel.Sigil,
	schlussel.Akte,
	schlussel.Bezeichnung,
	schlussel.Etikett,
	schlussel.Kennung,
	schlussel.Tai,
	schlussel.Typ,
	schlussel.MutterMetadateiMutterKennung,
	schlussel.ShaMetadateiMutterKennung,
	schlussel.ShaMetadatei,
	schlussel.VerzeichnisseEtikettImplicit,
	schlussel.VerzeichnisseEtikettExpanded,
	schlussel.VerzeichnisseEtiketten,
}

type Binary struct {
	bytes.Buffer
	BinaryField
	kennung.Sigil
	io.LimitedReader
}

//   ____                _
//  |  _ \ ___  __ _  __| |
//  | |_) / _ \/ _` |/ _` |
//  |  _ <  __/ (_| | (_| |
//  |_| \_\___|\__,_|\__,_|
//

func (bf *Binary) ReadFormatExactly(
	r io.ReaderAt,
	loc ennui.Loc,
	sk *Sku,
) (n int64, err error) {
	bf.BinaryField.Reset()
	bf.Buffer.Reset()

	var n1 int
	var n2 int64

	b := make([]byte, loc.ContentLength)

	n1, err = r.ReadAt(b, loc.Offset)
	n += int64(n1)

	if err == io.EOF {
		err = nil
	} else if err != nil {
		err = errors.Wrap(err)
		return
	}

	buf := bytes.NewBuffer(b)

	n1, err = ohio.ReadAllOrDieTrying(buf, bf.ContentLength[:])
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	_, _, err = bf.GetContentLength()

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = bf.readSigil(sk, buf)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	for buf.Len() > 0 {
		n2, err = bf.BinaryField.ReadFrom(buf)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = bf.readFieldKey(sk.Transacted); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

var ErrSkip = errors.New("skip")

type Sku struct {
	*sku.Transacted
	ennui.Range
	kennung.Sigil
}

func (bf *Binary) ReadFormatAndMatchSigil(
	r io.Reader,
	sk *Sku,
) (n int64, err error) {
	bf.BinaryField.Reset()
	bf.Buffer.Reset()

	var n1 int
	var n2 int64

	// loop thru entries to find the next one that matches the current sigil
	// when found, break the loop and deserialize it and return
	for {
		n1, err = ohio.ReadAllOrDieTrying(r, bf.ContentLength[:])
		n += int64(n1)

		if err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) && n == 0 {
				err = io.EOF
			}

			err = errors.WrapExcept(err, io.EOF)

			return
		}

		var contentLength64 int64
		_, contentLength64, err = bf.GetContentLength()

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		bf.R = r
		bf.N = contentLength64

		n2, err = bf.readSigil(sk, &bf.LimitedReader)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		if bf.Contains(sk.Sigil) {
			break
		}

		// TODO-P2 replace with buffered seeker
		// discard the next record
		if _, err = io.Copy(io.Discard, &bf.LimitedReader); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for bf.N > 0 {
		n2, err = bf.BinaryField.ReadFrom(&bf.LimitedReader)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = bf.readFieldKey(sk.Transacted); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

var errExpectedSigil = errors.New("expected sigil")

func (bf *Binary) readSigil(
	sk *Sku,
	r io.Reader,
) (n int64, err error) {
	n, err = bf.BinaryField.ReadFrom(r)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if bf.Schlussel != schlussel.Sigil {
		err = errors.Wrapf(errExpectedSigil, "Key: %s", bf.Schlussel)
		return
	}

	if _, err = sk.Sigil.ReadFrom(&bf.Content); err != nil {
		err = errors.Wrap(err)
		return
	}

	if sk.IncludesHidden() {
		sk.SetArchiviert(true)
	}

	return
}

func (bf *Binary) readFieldKey(
	sk *sku.Transacted,
) (err error) {
	switch bf.Schlussel {
	case schlussel.Akte:
		if _, err = sk.Metadatei.Akte.ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.Bezeichnung:
		if err = sk.Metadatei.Bezeichnung.Set(bf.Content.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.Etikett:
		var e kennung.Etikett

		if err = e.Set(bf.Content.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = sk.Metadatei.AddEtikettPtr(&e); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.Kennung:
		if _, err = sk.Kennung.ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.Tai:
		if _, err = sk.Metadatei.Tai.ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.Typ:
		if err = sk.Metadatei.Typ.Set(bf.Content.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.MutterMetadateiMutterKennung:
		if _, err = sk.Metadatei.Mutter().ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.ShaMetadateiMutterKennung:
		if _, err = sk.Metadatei.Sha().ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.ShaMetadatei:
		if _, err = sk.Metadatei.SelbstMetadatei.ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.VerzeichnisseEtikettImplicit:
		var e kennung.Etikett

		if err = e.Set(bf.Content.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = sk.Metadatei.Verzeichnisse.AddEtikettImplicitPtr(&e); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.VerzeichnisseEtikettExpanded:
		var e kennung.Etikett

		if err = e.Set(bf.Content.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = sk.Metadatei.Verzeichnisse.AddEtikettExpandedPtr(&e); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.VerzeichnisseEtiketten:
		var e etiketten_path.Path

		if _, err = e.ReadFrom(&bf.Content); err != nil {
			err = errors.WrapExcept(err, io.EOF)
			return
		}

		sk.Metadatei.Verzeichnisse.AddPath(&e)

	default:
		// panic(fmt.Sprintf("unsupported key: %s", key))
	}

	return
}

//  __        __    _ _
//  \ \      / / __(_) |_ ___
//   \ \ /\ / / '__| | __/ _ \
//    \ V  V /| |  | | ||  __/
//     \_/\_/ |_|  |_|\__\___|
//

func (bf *Binary) WriteFormat(
	w io.Writer,
	sk *sku.Transacted,
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

func (bf *Binary) writeFieldKey(
	sk *sku.Transacted,
) (n int64, err error) {
	switch bf.Schlussel {
	case schlussel.Sigil:
		s := bf.Sigil

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

func (bf *Binary) writeSha(
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

func (bf *Binary) writeFieldWriterTo(
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

func (bf *Binary) writeFieldBinaryMarshaler(
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

func (bf *Binary) writeFieldByteReader(
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

//   _____ _      _     _
//  |  ___(_) ___| | __| |___
//  | |_  | |/ _ \ |/ _` / __|
//  |  _| | |  __/ | (_| \__ \
//  |_|   |_|\___|_|\__,_|___/
//

type BinaryField struct {
	schlussel.Schlussel
	ContentLength [2]uint8
	Content       bytes.Buffer
}

func (bf *BinaryField) String() string {
	cl, _, _ := bf.GetContentLength()
	return fmt.Sprintf("%s:%d:%x", bf.Schlussel, cl, bf.Content.Bytes())
}

func (bf *BinaryField) Reset() {
	bf.Schlussel.Reset()
	bf.ContentLength[0] = 0
	bf.ContentLength[1] = 0
	bf.Content.Reset()
}

func (bf *BinaryField) GetContentLength() (contentLength int, contentLength64 int64, err error) {
	var n int
	contentLength64, n = binary.Varint(bf.ContentLength[:])

	if n <= 0 {
		err = errors.Errorf("error in content length: %d", n)
		return
	}

	if contentLength64 > math.MaxUint16 {
		err = errContentLengthTooLarge
		return
	}

	if contentLength64 < 0 {
		err = errContentLengthNegative
		return
	}

	return int(contentLength64), contentLength64, nil
}

func (bf *BinaryField) SetContentLength(v int) {
	if v < 0 {
		panic(errContentLengthNegative)
	}

	if v > math.MaxUint16 {
		panic(errContentLengthTooLarge)
	}

	binary.PutVarint(bf.ContentLength[:], int64(v))
}

var (
	errContentLengthTooLarge = errors.New("content length too large")
	errContentLengthNegative = errors.New("content length negative")
)

func (bf *BinaryField) ReadFrom(r io.Reader) (n int64, err error) {
	var n1 int
	var n2 int64
	n2, err = bf.Schlussel.ReadFrom(r)
	n += int64(n2)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	n1, err = ohio.ReadAllOrDieTrying(r, bf.ContentLength[:])
	n += int64(n1)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	contentLength, contentLength64, err := bf.GetContentLength()
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	bf.Content.Grow(contentLength)
	bf.Content.Reset()

	n2, err = io.CopyN(&bf.Content, r, contentLength64)
	n += n2

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	return
}

var errContentLengthDoesNotMatchContent = errors.New(
	"content length does not match content",
)

func (bf *BinaryField) WriteTo(w io.Writer) (n int64, err error) {
	if bf.Content.Len() > math.MaxUint16 {
		err = errContentLengthTooLarge
		return
	}

	bf.SetContentLength(bf.Content.Len())

	var n1 int
	var n2 int64
	n2, err = bf.Schlussel.WriteTo(w)
	n += int64(n2)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	n1, err = ohio.WriteAllOrDieTrying(w, bf.ContentLength[:])
	n += int64(n1)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	n2, err = io.Copy(w, &bf.Content)
	n += n2

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	return
}
