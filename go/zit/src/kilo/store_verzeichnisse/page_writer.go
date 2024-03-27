package store_verzeichnisse

import (
	"bufio"
	"os"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/files"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/golf/ennui"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type ShaTuple struct {
	Sha, Mutter *sha.Sha
	ennui.Range
	kennung.Sigil
}

type KennungShaMap map[string]ShaTuple

type pageWriter struct {
	*PageTuple
	binaryDecoder
	binaryEncoder
	*os.File
	bufio.Reader
	bufio.Writer

	ennui.Range
	offsetLast, offset int64
	kennungShaMap      KennungShaMap
}

func (pw *pageWriter) Flush() (err error) {
	if !pw.hasChanges {
		log.Log().Print("not flushing, no changes")
		return
	}

	defer pw.added.Reset()
	defer pw.addedSchwanz.Reset()

	pw.kennungShaMap = make(KennungShaMap)
	pw.binaryDecoder = makeBinary(kennung.SigilHistory)
	pw.binaryDecoder.Sigil = kennung.SigilHistory

	path := pw.Path()

	// If the cache file does not exist and we have nothing to add, short
	// circuit the flush. This condition occurs on the initial init when the
	// konfig is changed but there are no zettels yet.
	if !files.Exists(path) && pw.waitingToAddLen() == 0 {
		return
	}

	log.Log().Print("changesAreHistorical", pw.changesAreHistorical)
	log.Log().Print("added", pw.added.Len())
	log.Log().Print("addedSchwanz", pw.addedSchwanz.Len())

	if pw.added.Len() == 0 && !pw.changesAreHistorical {
		if pw.File, err = files.OpenReadWrite(path); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, pw.File)

		pw.Reader.Reset(pw.File)
		pw.Writer.Reset(pw.File)

		return pw.flushJustSchwanz()
	} else {
		if pw.File, err = pw.standort.FileTempLocal(); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloseAndRename(&err, pw.File, pw.Name(), path)

		pw.Reader.Reset(pw.File)
		pw.Writer.Reset(pw.File)

		return pw.flushBoth()
	}
}

func (pw *pageWriter) flushBoth() (err error) {
	chain := iter.MakeChain(
		pw.konfig.ApplyToSku,
		pw.writeOne,
	)

	if err = pw.CopyJustHistoryAndAdded(
		makeSigil(kennung.SigilHistory, kennung.SigilHidden),
		chain,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	for {
		popped, ok := pw.addedSchwanz.Pop()

		if !ok {
			break
		}

		if err = chain(popped); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = pw.Writer.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, st := range pw.kennungShaMap {
		st.Add(kennung.SigilSchwanzen)

		if err = pw.updateSigil(pw, st.Sigil, st.Offset); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (pw *pageWriter) flushJustSchwanz() (err error) {
	if err = pw.CopyJustHistoryFrom(
		&pw.Reader,
		makeSigil(kennung.SigilHistory, kennung.SigilHidden),
		func(sk Sku) (err error) {
			pw.Range = sk.Range
			pw.SaveSha(sk.Transacted, sk.Sigil)
			return
		},
	); err != nil {
		err = errors.Wrapf(err, "Page: %s", pw.PageId)
		return
	}

	chain := iter.MakeChain(
		pw.konfig.ApplyToSku,
		pw.removeOldSchwanzen,
		pw.writeOne,
	)

	for {
		popped, ok := pw.addedSchwanz.Pop()

		if !ok {
			break
		}

		if err = chain(popped); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = pw.Writer.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, st := range pw.kennungShaMap {
		st.Add(kennung.SigilSchwanzen)

		if err = pw.updateSigil(pw, st.Sigil, st.Offset); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (pw *pageWriter) writeOne(
	z *sku.Transacted,
) (err error) {
	pw.Offset += pw.ContentLength

	if pw.ContentLength, err = pw.writeFormat(
		&pw.Writer,
		skuWithSigil{Transacted: z},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	pw.etikettIndex.Add(z.Metadatei.GetEtiketten())

	pw.SaveSha(z, kennung.SigilHistory)

	return
}

func (pw *pageWriter) SaveSha(z *sku.Transacted, sigil kennung.Sigil) {
	k := z.GetKennung()

	record := pw.kennungShaMap[k.String()]
	record.Range = pw.Range

	if z.Metadatei.Verzeichnisse.Archiviert.Bool() {
		sigil.Add(kennung.SigilHidden)
	}

	record.Sigil = sigil

	if z.Metadatei.Verzeichnisse.Archiviert.Bool() {
		record.Add(kennung.SigilHidden)
	} else {
		record.Del(kennung.SigilHidden)
	}

	pw.kennungShaMap[k.String()] = record
}

func (pw *pageWriter) removeOldSchwanzen(sk *sku.Transacted) (err error) {
	ks := sk.Kennung.String()
	st, ok := pw.kennungShaMap[ks]

	if !ok {
		return
	}

	st.Del(kennung.SigilSchwanzen)

	if err = pw.updateSigil(pw, st.Sigil, st.Offset); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
