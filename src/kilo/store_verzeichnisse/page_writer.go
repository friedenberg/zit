package store_verzeichnisse

import (
	"bufio"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/ennui"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/sku_fmt"
)

type ShaTuple struct {
	Sha, Mutter           *sha.Sha
	Offset, ContentLength int64
	kennung.Sigil
}

type KennungShaMap map[string]ShaTuple

type pageWriter struct {
	*PageTuple
	sku_fmt.Binary
	*os.File
	bufio.Writer

	offsetLast, offset int64
	kennungShaMap      KennungShaMap
}

func (pw *pageWriter) ReadMutter(z *sku.Transacted) (err error) {
	k := z.GetKennung()
	old := pw.kennungShaMap[k.String()]

	if !old.Mutter.IsNull() {
		if err = z.GetMetadatei().Mutter.SetShaLike(old.Mutter); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = z.CalculateObjekteSha(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (pw *pageWriter) writeOne(
	z *sku.Transacted,
) (err error) {
	pw.offsetLast = pw.offset

	var n int64

	if n, err = pw.WriteFormat(&pw.Writer, z); err != nil {
		err = errors.Wrap(err)
		return
	}

	pw.etikettIndex.Add(z.Metadatei.GetEtiketten())

	if pw.ennuiShas == nil {
		return
	}

	if err = pw.ennuiShas.AddMetadatei(
		z.GetMetadatei(),
		ennui.Loc{
			Page:          pw.Index,
			Offset:        pw.offsetLast,
			ContentLength: n,
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	pw.offset += int64(n)

	return
}

func (pw *pageWriter) SaveSha(z *sku.Transacted) (err error) {
	k := z.GetKennung()
	var sh sha.Sha

	if err = sh.SetShaLike(&z.GetMetadatei().Sha); err != nil {
		err = errors.Wrap(err)
		return
	}

	old := pw.kennungShaMap[k.String()]
	old.Mutter = old.Sha
	old.Sha = &sh
	old.Offset = pw.offsetLast
	old.ContentLength = pw.offset - pw.offsetLast
	old.Sigil = kennung.SigilHistory

	if z.Metadatei.Verzeichnisse.Archiviert.Bool() {
		old.Add(kennung.SigilHidden)
	}

	pw.kennungShaMap[k.String()] = old

	return
}

func (pw *pageWriter) flush() (err error) {
	if !pw.hasChanges {
		return
	}

	defer pw.added.Reset()
	defer pw.addedSchwanz.Reset()

	pw.kennungShaMap = make(KennungShaMap)
	pw.Binary = sku_fmt.Binary{Sigil: kennung.SigilHistory}

	path := pw.Path()

	// If the cache file does not exist and we have nothing to add, short circuit
	// the flush. This condition occurs on the initial init when the konfig is
	// changed but there are no zettels yet.
	if !files.Exists(path) && pw.waitingToAddLen() == 0 {
		return
	}

	if pw.File, err = pw.standort.FileTempLocal(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloseAndRename(&err, pw.File, pw.Name(), path)

	pw.Reset(pw.File)

	chain := iter.MakeChain(
		pw.konfig.ApplyToSku,
		pw.ReadMutter,
		pw.writeOne,
		pw.SaveSha,
	)

	if err = pw.Copy(kennung.SigilHistory, chain); err != nil {
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

	var n int

	for ks, st := range pw.kennungShaMap {
		st.Add(kennung.SigilSchwanzen)

		// 2 uint8 + offset + 2 uint8 + Schlussel
		offset := int64(2) + st.Offset + int64(3)

		if n, err = pw.WriteAt([]byte{st.Byte()}, offset); err != nil {
			err = errors.Wrap(err)
			return
		}

		if n != 1 {
			panic(errors.Errorf("expected 1 byte but wrote %d", n))
		}

		shK := sha.FromString(ks)

		if err = pw.ennuiKennung.AddSha(
			shK,
			ennui.Loc{
				Page:          pw.Index,
				Offset:        st.Offset,
				ContentLength: st.ContentLength,
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
