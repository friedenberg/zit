package stream_index

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
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

var binaryFieldOrder = []keys.Binary{
	keys.Sigil,
	keys.ObjectId,
	keys.Blob,
	keys.Description,
	keys.Tag,
	keys.Tai,
	keys.Type,
	keys.ShaParentMetadataParentObjectId,
	keys.ShaMetadataParentObjectId,
	keys.ShaMetadata,
	keys.ShaMetadataWithoutTai,
	keys.CacheParentTai,
	keys.CacheTagImplicit,
	keys.CacheTagExpanded,
	keys.CacheTags,
}

func makeBinary(s ids.Sigil) binaryDecoder {
	return binaryDecoder{
		PrimitiveQueryGroup: sku.MakePrimitiveQueryGroupWithSigils(s),
		Sigil:               s,
	}
}

func makeBinaryWithQueryGroup(
	qg sku.PrimitiveQueryGroup,
	s ids.Sigil,
) binaryDecoder {
	ui.Log().Print(qg)
	if qg == nil {
		qg = sku.MakePrimitiveQueryGroup()
	}

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
	sk *skuWithRangeAndSigil,
) (n int64, err error) {
	bf.binaryField.Reset()
	bf.Buffer.Reset()

	var n1 int
	var n2 int64

	b := make([]byte, sk.ContentLength)

	n1, err = r.ReadAt(b, sk.Offset)
	n += int64(n1)

	if err == io.EOF {
		err = nil
	} else if err != nil {
		err = errors.Wrap(err)
		return
	}

	buf := bytes.NewBuffer(b)

	n1, bf.ContentLength, err = ohio.ReadFixedUInt16(buf)
	n += int64(n1)

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
		n1, bf.ContentLength, err = ohio.ReadFixedUInt16(r)
		n += int64(n1)

		if err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) && n == 0 {
				err = io.EOF
			}

			err = errors.WrapExcept(err, io.EOF)

			return
		}

		contentLength64 := int64(bf.ContentLength)

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

		g := genres.Must(sk.Transacted)
		q, ok := bf.Get(g)

		// TODO-D4 use query to decide whether to read and inflate or skip
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

			if q.ContainsObjectId(&sk.ObjectId) &&
				(qs.ContainsOneOf(ids.SigilHistory) ||
					sk.ContainsOneOf(ids.SigilLatest)) {
				break
			}
		}

		// TODO-D4 replace with buffered seeker
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

	if bf.Binary != keys.Sigil {
		err = errors.Wrapf(errExpectedSigil, "Key: %s", bf.Binary)
		return
	}

	if _, err = sk.Sigil.ReadFrom(&bf.Content); err != nil {
		err = errors.Wrap(err)
		return
	}

	sk.SetDormant(sk.IncludesHidden())

	return
}

func (bf *binaryDecoder) readFieldKey(
	sk *sku.Transacted,
) (err error) {
	switch bf.Binary {
	case keys.Blob:
		if _, err = sk.Metadata.Blob.ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.Description:
		if err = sk.Metadata.Description.Set(bf.Content.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.Tag:
		var e ids.Tag

		if err = e.Set(bf.Content.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = sk.AddTagPtrFast(&e); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.ObjectId:
		if _, err = sk.ObjectId.ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.Tai:
		if _, err = sk.Metadata.Tai.ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.CacheParentTai:
		if _, err = sk.Metadata.Cache.ParentTai.ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.Type:
		if err = sk.Metadata.Type.Set(bf.Content.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.ShaParentMetadataParentObjectId:
		if _, err = sk.Metadata.Mutter().ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.ShaMetadataParentObjectId:
		if _, err = sk.Metadata.Sha().ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.ShaMetadata:
		if _, err = sk.Metadata.SelfMetadata.ReadFrom(&bf.Content); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.ShaMetadataWithoutTai:
		if _, err = sk.Metadata.SelfMetadataWithoutTai.ReadFrom(
			&bf.Content,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.CacheTagImplicit:
		var e ids.Tag

		if err = e.Set(bf.Content.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = sk.Metadata.Cache.AddTagsImplicitPtr(&e); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.CacheTagExpanded:
		var e ids.Tag

		if err = e.Set(bf.Content.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = sk.Metadata.Cache.AddTagExpandedPtr(&e); err != nil {
			err = errors.Wrap(err)
			return
		}

	case keys.CacheTags:
		var e tag_paths.PathWithType

		if _, err = e.ReadFrom(&bf.Content); err != nil {
			err = errors.WrapExcept(err, io.EOF)
			return
		}

		sk.Metadata.Cache.TagPaths.AddPath(&e)

	default:
		// panic(fmt.Sprintf("unsupported key: %s", key))
	}

	return
}
