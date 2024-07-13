package store_verzeichnisse

import (
	"bytes"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/keys"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/tag_paths"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_probe_index"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

var binaryFieldOrder = []keys.Key{
	keys.Sigil,
	keys.ObjectId,
	keys.Blob,
	keys.Description,
	keys.Tag,
	keys.Tai,
	keys.Type,
	keys.MutterMetadateiMutterKennung,
	keys.ShaMetadateiMutterKennung,
	keys.ShaMetadatei,
	keys.ShaMetadateiSansTai,
	keys.VerzeichnisseEtikettImplicit,
	keys.VerzeichnisseEtikettExpanded,
	keys.VerzeichnisseEtiketten,
}

func makeFlushQueryGroup(ss ...ids.Sigil) sku.PrimitiveQueryGroup {
	return &flushQueryGroup{Sigil: ids.MakeSigil(ss...)}
}

func makeBinary(s ids.Sigil) binaryDecoder {
	return binaryDecoder{
		PrimitiveQueryGroup: makeFlushQueryGroup(s),
		Sigil:               s,
	}
}

func makeBinaryWithQueryGroup(
	qg sku.PrimitiveQueryGroup,
	s ids.Sigil,
) binaryDecoder {
	ui.Log().Print(qg)
	if !qg.HasHidden() {
		s.Add(ids.SigilHidden)
	}

	return binaryDecoder{
		PrimitiveQueryGroup: qg,
		Sigil:               s,
	}
}

type binaryDecoder struct {
	bytes.Buffer
	binaryField
	ids.Sigil
	sku.PrimitiveQueryGroup
	io.LimitedReader
}

func (bf *binaryDecoder) readFormatExactly(
	r io.ReaderAt,
	loc object_probe_index.Loc,
	sk *skuWithRangeAndSigil,
) (n int64, err error) {
	bf.binaryField.Reset()
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
		n2, err = bf.binaryField.ReadFrom(buf)
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

func (bf *binaryDecoder) readFormatAndMatchSigil(
	r io.Reader,
	sk *skuWithRangeAndSigil,
) (n int64, err error) {
	bf.binaryField.Reset()
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

		n2, err = bf.binaryField.ReadFrom(&bf.LimitedReader)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = bf.readFieldKey(sk.Transacted); err != nil {
			err = errors.Wrap(err)
			return
		}

		q, ok := bf.Get(genres.Must(sk.Transacted))

		// TODO use query to decide whether to read and inflate or skip
		if ok {
			qs := q.GetSigil()

			wantsHidden := qs.IncludesHidden()
			wantsHistory := qs.IncludesHistory()
			isSchwanzen := sk.Contains(ids.SigilLatest)
			isHidden := sk.Contains(ids.SigilHidden)

			// log.Log().Print(sk)
			// log.Log().Print("wantsHistory", wantsHistory)
			// log.Log().Print("wantsHidden", wantsHidden)
			// log.Log().Print("isSchwanzen", isSchwanzen)
			// log.Log().Print("isHidden", isHidden)

			if (wantsHistory && wantsHidden) ||
				(wantsHidden && isSchwanzen) ||
				(wantsHistory && !isHidden) ||
				(isSchwanzen && !isHidden) {
				break
			}

			if q.ContainsKennung(&sk.Kennung) &&
				(qs.ContainsOneOf(ids.SigilHistory) ||
					sk.ContainsOneOf(ids.SigilLatest)) {
				break
			}
		}

		// TODO-P2 replace with buffered seeker
		// discard the next record
		if _, err = io.Copy(io.Discard, &bf.LimitedReader); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for bf.N > 0 {
		n2, err = bf.binaryField.ReadFrom(&bf.LimitedReader)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = bf.readFieldKey(sk.Transacted); err != nil {
			err = errors.Wrapf(err, "Sku: %#v", sk.Transacted)
			return
		}
	}

	return
}

var errExpectedSigil = errors.New("expected sigil")

func (bf *binaryDecoder) readSigil(
	sk *skuWithRangeAndSigil,
	r io.Reader,
) (n int64, err error) {
	n, err = bf.binaryField.ReadFrom(r)
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if bf.Key != keys.Sigil {
		err = errors.Wrapf(errExpectedSigil, "Key: %s", bf.Key)
		return
	}

	if _, err = sk.Sigil.ReadFrom(&bf.Content); err != nil {
		err = errors.Wrap(err)
		return
	}

	sk.SetSchlummernd(sk.IncludesHidden())

	return
}

func (bf *binaryDecoder) readFieldKey(
	sk *sku.Transacted,
) (err error) {
	switch bf.Key {
	case keys.Blob:
		if _, err = sk.Metadatei.Akte.ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.Description:
		if err = sk.Metadatei.Description.Set(bf.Content.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.Tag:
		var e ids.Tag

		if err = e.Set(bf.Content.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = sk.AddEtikettPtrFast(&e); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.ObjectId:
		if _, err = sk.Kennung.ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.Tai:
		if _, err = sk.Metadatei.Tai.ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.Type:
		if err = sk.Metadatei.Type.Set(bf.Content.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.MutterMetadateiMutterKennung:
		if _, err = sk.Metadatei.Mutter().ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.ShaMetadateiMutterKennung:
		if _, err = sk.Metadatei.Sha().ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.ShaMetadatei:
		if _, err = sk.Metadatei.SelbstMetadatei.ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.ShaMetadateiSansTai:
		if _, err = sk.Metadatei.SelbstMetadateiSansTai.ReadFrom(
			&bf.Content,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.VerzeichnisseEtikettImplicit:
		var e ids.Tag

		if err = e.Set(bf.Content.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = sk.Metadatei.Cached.AddEtikettImplicitPtr(&e); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.VerzeichnisseEtikettExpanded:
		var e ids.Tag

		if err = e.Set(bf.Content.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = sk.Metadatei.Cached.AddEtikettExpandedPtr(&e); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.VerzeichnisseEtiketten:
		var e tag_paths.PathWithType

		if _, err = e.ReadFrom(&bf.Content); err != nil {
			err = errors.WrapExcept(err, io.EOF)
			return
		}

		sk.Metadatei.Cached.Etiketten.AddPath(&e)

	default:
		// panic(fmt.Sprintf("unsupported key: %s", key))
	}

	return
}
