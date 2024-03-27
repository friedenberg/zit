package sku_fmt

import (
	"bytes"
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/delta/ohio"
	"code.linenisgreat.com/zit/src/delta/schlussel"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/etiketten_path"
	"code.linenisgreat.com/zit/src/golf/ennui"
	"code.linenisgreat.com/zit/src/hotel/sku"
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
	schlussel.VerzeichnisseEtikettImplicit,
	schlussel.VerzeichnisseEtikettExpanded,
	schlussel.VerzeichnisseEtiketten,
}

func MakeSigil(ss ...kennung.Sigil) sku.MatcherGroup {
	return &NopSigil{Sigil: kennung.MakeSigil(ss...)}
}

type NopSigil struct {
	kennung.Sigil
}

func (qg *NopSigil) Get(_ gattung.Gattung) (sku.Query, bool) {
	return qg, true
}

func (s *NopSigil) ContainsMatchable(_ *sku.Transacted) bool {
	return true
}

func (s *NopSigil) String() string {
	panic("should never be called")
}

func (s *NopSigil) ContainsKennung(_ *kennung.Kennung2) bool {
	return false
}

func (s *NopSigil) GetSigil() kennung.Sigil {
	return s.Sigil
}

func (s *NopSigil) Each(_ schnittstellen.FuncIter[sku.QueryBase]) error {
	return nil
}

func MakeBinary(s kennung.Sigil) Binary {
	return Binary{
		MatcherGroup: MakeSigil(s),
		Sigil:        s,
	}
}

func MakeBinaryWithQueryGroup(qg sku.MatcherGroup, s kennung.Sigil) Binary {
	return Binary{
		MatcherGroup: qg,
		Sigil:        s,
	}
}

type Binary struct {
	bytes.Buffer
	BinaryField
	kennung.Sigil
	sku.MatcherGroup
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

		q, ok := bf.Get(gattung.Must(sk.Transacted))

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
		n2, err = bf.BinaryField.ReadFrom(&bf.LimitedReader)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = bf.readFieldKey(sk.Transacted); err != nil {
			err = errors.Wrapf(err, "Sku: %#v", sk)
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
