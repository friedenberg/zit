package objekte_store

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

type StoredParseSaver[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
] interface {
	ParseSaveStored(
		sem sku.ExternalMaybe,
		t *sku.External,
	) (a OPtr, err error)
}

type storedParserSaver[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
] struct {
	awf        schnittstellen.AkteWriterFactory
	akteParser objekte.AkteParser[OPtr]
}

func MakeStoredParseSaver[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
](
	owf schnittstellen.ObjekteIOFactory,
	awf schnittstellen.AkteIOFactory,
	akteParser objekte.AkteParser[OPtr],
	pmf objekte_format.Format,
	op objekte_format.Options,
) storedParserSaver[O, OPtr] {
	if akteParser == nil {
		akteParser = objekte.MakeNopAkteParseSaver[O, OPtr](awf)
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
	sem sku.ExternalMaybe,
	t *sku.External,
) (o OPtr, err error) {
	var f *os.File

	errors.TodoP2("support akte")
	if f, err = files.OpenExclusiveReadOnly(sem.FDs.Objekte.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.FDs = sem.FDs

	if err = t.Kennung.SetWithKennung(sem.Kennung); err != nil {
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

	errors.Log().Printf("parsed %d akte bytes", n)

	return
}
