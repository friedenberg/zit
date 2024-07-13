package store_verzeichnisse

import (
	"bytes"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/schlussel"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/etiketten_path"
	"code.linenisgreat.com/zit/go/zit/src/golf/ennui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

var binaryFieldOrder = []schlussel.Schlussel{
	schlussel.Sigil,
	schlussel.Kennung,
	schlussel.Akte,
	schlussel.Bezeichnung,
	schlussel.Etikett,
	schlussel.Tai,
	schlussel.Typ,
	schlussel.MutterMetadateiMutterKennung,
	schlussel.ShaMetadateiMutterKennung,
	schlussel.ShaMetadatei,
	schlussel.ShaMetadateiSansTai,
	schlussel.VerzeichnisseEtikettImplicit,
	schlussel.VerzeichnisseEtikettExpanded,
	schlussel.VerzeichnisseEtiketten,
}

func makeFlushQueryGroup(ss ...kennung.Sigil) sku.PrimitiveQueryGroup {
	return &flushQueryGroup{Sigil: kennung.MakeSigil(ss...)}
}

func makeBinary(s kennung.Sigil) binaryDecoder {
	return binaryDecoder{
		PrimitiveQueryGroup: makeFlushQueryGroup(s),
		Sigil:               s,
	}
}

func makeBinaryWithQueryGroup(
	qg sku.PrimitiveQueryGroup,
	s kennung.Sigil,
) binaryDecoder {
	ui.Log().Print(qg)
	if !qg.HasHidden() {
		s.Add(kennung.SigilHidden)
	}

	return binaryDecoder{
		PrimitiveQueryGroup: qg,
		Sigil:               s,
	}
}

type binaryDecoder struct {
	bytes.Buffer
	binaryField
	kennung.Sigil
	sku.PrimitiveQueryGroup
	io.LimitedReader
}

func (bf *binaryDecoder) readFormatExactly(
	r io.ReaderAt,
	loc ennui.Loc,
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

		q, ok := bf.Get(gattung.Must(sk.Transacted))

		// TODO use query to decide whether to read and inflate or skip
		if ok {
			qs := q.GetSigil()

			wantsHidden := qs.IncludesHidden()
			wantsHistory := qs.IncludesHistory()
			isSchwanzen := sk.Contains(kennung.SigilSchwanzen)
			isHidden := sk.Contains(kennung.SigilHidden)

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
				(qs.ContainsOneOf(kennung.SigilHistory) ||
					sk.ContainsOneOf(kennung.SigilSchwanzen)) {
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

	if bf.Schlussel != schlussel.Sigil {
		err = errors.Wrapf(errExpectedSigil, "Key: %s", bf.Schlussel)
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
		var e kennung.Tag

		if err = e.Set(bf.Content.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = sk.AddEtikettPtrFast(&e); err != nil {
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

	case schlussel.ShaMetadateiSansTai:
		if _, err = sk.Metadatei.SelbstMetadateiSansTai.ReadFrom(
			&bf.Content,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.VerzeichnisseEtikettImplicit:
		var e kennung.Tag

		if err = e.Set(bf.Content.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = sk.Metadatei.Verzeichnisse.AddEtikettImplicitPtr(&e); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.VerzeichnisseEtikettExpanded:
		var e kennung.Tag

		if err = e.Set(bf.Content.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = sk.Metadatei.Verzeichnisse.AddEtikettExpandedPtr(&e); err != nil {
			err = errors.Wrap(err)
			return
		}

	case schlussel.VerzeichnisseEtiketten:
		var e etiketten_path.PathWithType

		if _, err = e.ReadFrom(&bf.Content); err != nil {
			err = errors.WrapExcept(err, io.EOF)
			return
		}

		sk.Metadatei.Verzeichnisse.Etiketten.AddPath(&e)

	default:
		// panic(fmt.Sprintf("unsupported key: %s", key))
	}

	return
}
