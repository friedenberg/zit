package store_verzeichnisse

import (
	"bufio"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_probe_index"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// type entry struct {
// 	Sha, Mutter *sha.Sha
// 	ennui.Range
// 	kennung.Sigil
// }

type KennungShaMap map[string]skuWithRangeAndSigil

type writer struct {
	*Page
	binaryDecoder
	binaryEncoder
	*os.File
	bufio.Reader
	bufio.Writer

	changesAreHistorical bool

	object_probe_index.Range
	offsetLast, offset int64
	kennungShaMap      KennungShaMap
}

func (pw *writer) Flush() (err error) {
	if !pw.hasChanges {
		ui.Log().Print("not flushing, no changes")
		return
	}

	defer pw.added.Reset()
	defer pw.addedSchwanz.Reset()

	pw.kennungShaMap = make(KennungShaMap)
	pw.binaryDecoder = makeBinary(ids.SigilHistory)
	pw.binaryDecoder.Sigil = ids.SigilHistory

	path := pw.Path()

	// If the cache file does not exist and we have nothing to add, short
	// circuit the flush. This condition occurs on the initial init when the
	// konfig is changed but there are no zettels yet.
	if !files.Exists(path) && pw.waitingToAddLen() == 0 {
		return
	}

	ui.Log().Print("changesAreHistorical", pw.changesAreHistorical)
	ui.Log().Print("added", pw.added.Len())
	ui.Log().Print("addedSchwanz", pw.addedSchwanz.Len())

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
		if pw.File, err = pw.fs_home.FileTempLocal(); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloseAndRename(&err, pw.File, pw.Name(), path)

		pw.Reader.Reset(pw.File)
		pw.Writer.Reset(pw.File)

		return pw.flushBoth()
	}
}

func (pw *writer) flushBoth() (err error) {
	ui.Log().Printf("flushing both: %s", pw.Path())

	chain := iter.MakeChain(
		pw.konfig.ApplySchlummerndAndRealizeEtiketten,
		pw.writeOne,
	)

	if err = pw.copyJustHistoryAndAdded(
		makeFlushQueryGroup(ids.SigilHistory, ids.SigilHidden),
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
		if err = pw.updateSigilWithSchwanzen(st); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (pw *writer) updateSigilWithSchwanzen(st skuWithRangeAndSigil) (err error) {
	st.Add(ids.SigilLatest)

	if err = pw.WriteOneObjekteMetadatei(st.Transacted); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = pw.updateSigil(pw, st.Sigil, st.Offset); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (pw *writer) flushJustSchwanz() (err error) {
	ui.Log().Printf("flushing just schwanz: %s", pw.Path())

	if err = pw.copyJustHistoryFrom(
		&pw.Reader,
		makeFlushQueryGroup(ids.SigilHistory, ids.SigilHidden),
		func(sk skuWithRangeAndSigil) (err error) {
			pw.Range = sk.Range
			pw.saveSchwanz(sk.Transacted, sk.Sigil)
			return
		},
	); err != nil {
		err = errors.Wrapf(err, "Page: %s", pw.PageId)
		return
	}

	chain := iter.MakeChain(
		pw.konfig.ApplySchlummerndAndRealizeEtiketten,
		pw.removeOldSchwanz,
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
		if err = pw.updateSigilWithSchwanzen(st); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (pw *writer) writeOne(
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

	pw.saveSchwanz(z, ids.SigilHistory)

	return
}

func (pw *writer) saveSchwanz(z *sku.Transacted, sigil ids.Sigil) {
	k := z.GetKennung()
	ks := k.String()

	record := pw.kennungShaMap[ks]
	record.Range = pw.Range

	if record.Transacted == nil {
		record.Transacted = sku.GetTransactedPool().Get()
	}

	sku.TransactedResetter.ResetWith(record.Transacted, z)

	record.Sigil = sigil

	if z.Metadatei.Verzeichnisse.Schlummernd.Bool() {
		record.Add(ids.SigilHidden)
	} else {
		record.Del(ids.SigilHidden)
	}

	pw.kennungShaMap[ks] = record
}

func (pw *writer) removeOldSchwanz(sk *sku.Transacted) (err error) {
	ks := sk.Kennung.String()
	st, ok := pw.kennungShaMap[ks]

	if !ok {
		return
	}

	st.Del(ids.SigilLatest)

	if err = pw.updateSigil(pw, st.Sigil, st.Offset); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
