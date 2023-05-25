package objekte_store

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/objekte_format"
)

type StoredParseSaver[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
] interface {
	ParseSaveStored(
		sem sku.ExternalMaybe[K, KPtr],
		t *objekte.External[O, OPtr, K, KPtr],
	) (err error)
}

type storedParserSaver[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
] struct {
	awf          schnittstellen.AkteWriterFactory
	akteParser   objekte.AkteParseSaver[OPtr]
	objekteSaver ObjekteSaver
}

func MakeStoredParseSaver[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
](
	owf schnittstellen.ObjekteIOFactory,
	awf schnittstellen.AkteIOFactory,
	akteParser objekte.AkteParseSaver[OPtr],
	pmf objekte_format.Format,
) storedParserSaver[O, OPtr, K, KPtr] {
	if akteParser == nil {
		akteParser = objekte.MakeNopAkteParseSaver[O, OPtr](awf)
	}

	if pmf == nil {
		panic("persisted_metadatei_format.Format was nil")
	}

	return storedParserSaver[O, OPtr, K, KPtr]{
		awf:        awf,
		akteParser: akteParser,
		objekteSaver: MakeObjekteSaver(
			owf,
			pmf,
		),
	}
}

func (h storedParserSaver[O, OPtr, K, KPtr]) ParseSaveStored(
	sem sku.ExternalMaybe[K, KPtr],
	t *objekte.External[O, OPtr, K, KPtr],
) (err error) {
	var f *os.File

	errors.TodoP2("support akte")
	if f, err = files.OpenExclusiveReadOnly(sem.FDs.Objekte.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.Sku.FDs = sem.FDs
	t.Sku.Kennung = sem.Kennung

	r := sha.MakeReadCloser(f)

	defer errors.DeferredCloser(&err, r)

	var akteSha schnittstellen.Sha

	if akteSha, _, err = h.readAkte(r, &t.Akte); err != nil {
		err = errors.Wrap(err)
		return
	}

	readerSha := sha.Make(r.Sha())

	if !readerSha.EqualsSha(akteSha) {
		err = errors.Errorf(
			"akte reader got %s but AkteParseSaver got %s",
			readerSha,
			akteSha,
		)

		return
	}

	t.SetAkteSha(readerSha)

	if err = h.objekteSaver.SaveObjekte(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (h storedParserSaver[O, OPtr, K, KPtr]) readAkte(
	r sha.ReadCloser,
	o OPtr,
) (sh schnittstellen.Sha, n int64, err error) {
	if sh, n, err = h.akteParser.ParseSaveAkte(r, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("parsed %d akte bytes", n)

	return
}
