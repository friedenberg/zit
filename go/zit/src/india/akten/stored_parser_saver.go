package akten

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/store_fs"
)

type StoredParseSaver[
	O schnittstellen.Akte[O],
	OPtr schnittstellen.AktePtr[O],
] interface {
	ParseSaveStored(
		sem sku.KennungFDPair,
		t *store_fs.External,
	) (a OPtr, err error)
}

type storedParserSaver[
	O schnittstellen.Akte[O],
	OPtr schnittstellen.AktePtr[O],
] struct {
	awf        schnittstellen.AkteWriterFactory
	akteParser Parser[O, OPtr]
}

func MakeStoredParseSaver[
	O schnittstellen.Akte[O],
	OPtr schnittstellen.AktePtr[O],
](
	awf schnittstellen.AkteIOFactory,
	akteParser Parser[O, OPtr],
	pmf objekte_format.Format,
	op objekte_format.Options,
) storedParserSaver[O, OPtr] {
	if akteParser == nil {
		akteParser = MakeNopAkteParseSaver[O, OPtr](awf)
	}

	if pmf == nil {
		panic("persisted_metadatei_format.Format was nil")
	}

	return storedParserSaver[O, OPtr]{
		awf:        awf,
		akteParser: akteParser,
	}
}

func (h storedParserSaver[O, OPtr]) ParseSaveStored(
	sem *sku.KennungFDPair,
	t *store_fs.External,
) (o OPtr, err error) {
	var f *os.File

	errors.TodoP2("support akte")
	if f, err = files.OpenExclusiveReadOnly(sem.FDs.Objekte.GetPath()); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.FDs = sem.FDs

	if err = t.Kennung.SetWithKennung(&sem.Kennung); err != nil {
		err = errors.Wrap(err)
		return
	}

	r := sha.MakeReadCloser(f)

	defer errors.DeferredCloser(&err, r)

	var akteSha schnittstellen.ShaLike

	// TODO-P3 switch to pool
	var o1 O
	o = OPtr(&o1)

	if akteSha, _, err = h.readAkte(r, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	readerSha := sha.Make(r.GetShaLike())

	if !readerSha.EqualsSha(akteSha) {
		err = errors.Errorf(
			"akte reader got %s but AkteParseSaver got %s",
			readerSha,
			akteSha,
		)

		return
	}

	t.SetAkteSha(readerSha)

	return
}

func (h storedParserSaver[O, OPtr]) readAkte(
	r sha.ReadCloser,
	o OPtr,
) (sh schnittstellen.ShaLike, n int64, err error) {
	sw := sha.MakeWriter(io.Discard)

	if n, err = h.akteParser.ParseAkte(io.TeeReader(r, sw), o); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = sw.GetShaLike()

	ui.Log().Printf("parsed %d akte bytes", n)

	return
}
