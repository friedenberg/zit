package bestandsaufnahme

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/golf/ennui_shas"
	"code.linenisgreat.com/zit/go/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type Ennui interface {
	WriteOneObjekteMetadatei(o *sku.Transacted) (err error)
	ReadOneEnnui(*sha.Sha) (*sku.Transacted, error)
	ReadOneKennung(kennung.Kennung) (*sku.Transacted, error)
	ReadOneKennungSha(kennung.Kennung) (*sha.Sha, error)
}

type ennuiStore struct {
	standort                  standort.Standort
	persistentMetadateiFormat objekte_format.Format
	ennuiKennung              ennui_shas.Ennui
	options                   objekte_format.Options
}

func (s *ennuiStore) Initialize(
	standort standort.Standort,
	persistentMetadateiFormat objekte_format.Format,
	options objekte_format.Options,
) (err error) {
	s.standort = standort
	s.persistentMetadateiFormat = persistentMetadateiFormat
	s.options = options

	if s.ennuiKennung, err = ennui_shas.MakeNoDuplicates(
		s.standort,
		s.standort.DirVerzeichnisseVerweise(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *ennuiStore) ReadOneEnnui(sh *sha.Sha) (sk *sku.Transacted, err error) {
	var r sha.ReadCloser

	if r, err = s.standort.AkteReaderFrom(
		sh,
		s.standort.DirVerzeichnisseMetadateiKennungMutter(),
	); err != nil {
		if errors.IsNotExist(err) {
			err = collections.MakeErrNotFound(sh)
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.DeferredCloser(&err, r)

	rb := catgut.MakeRingBuffer(r, 0)

	sk = sku.GetTransactedPool().Get()

	var n int64

	n, err = s.persistentMetadateiFormat.ParsePersistentMetadatei(
		rb,
		sk,
		s.options,
	)

	if err == io.EOF && n > 0 {
		err = nil
	} else if err != io.EOF && err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = sk.CalculateObjekteShas(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *ennuiStore) ReadOneKennungSha(k kennung.Kennung) (sh *sha.Sha, err error) {
	left := sha.FromString(k.String())
	defer sha.GetPool().Put(left)

	if sh, err = s.ennuiKennung.ReadOne(left); err != nil {
		err = errors.Wrapf(err, "Kennung: %q, Left: %s", k, left)
		return
	}

	return
}

func (s *ennuiStore) ReadOneKennung(k kennung.Kennung) (sk *sku.Transacted, err error) {
	sh, err := s.ReadOneKennungSha(k)
	defer sha.GetPool().Put(sh)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if sk, err = s.ReadOneEnnui(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !sh.Equals(sk.Metadatei.Sha()) {
		err = errors.Errorf(
			"expected sha %q but got %q",
			sh,
			sk.Metadatei.Sha(),
		)

		return
	}

	return
}

func (s *ennuiStore) makeWriteMetadateiFunc(
	dir string,
	fo objekte_format.FormatGeneric,
	o *sku.Transacted,
	expected *sha.Sha,
) schnittstellen.FuncError {
	return func() (err error) {
		var sw sha.WriteCloser

		if sw, err = s.standort.AkteWriterToLight(dir); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, sw)

		if _, err = fo.WriteMetadateiTo(sw, o); err != nil {
			err = errors.Wrap(err)
			return
		}

		actual := sw.GetShaLike()

		if !expected.EqualsSha(actual) {
			err = errors.Errorf(
				"expected %q but got %q",
				expected,
				actual,
			)

			return
		}

		return
	}
}

func (s *ennuiStore) MakeFuncSaveOneVerweise(o *sku.Transacted) func() error {
	return func() (err error) {
		k := o.GetKennung()
		sh := sha.FromString(k.String())
		defer sha.GetPool().Put(sh)

		if err = s.ennuiKennung.AddSha(sh, o.Metadatei.Sha()); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (s *ennuiStore) WriteOneObjekteMetadatei(o *sku.Transacted) (err error) {
	if o.Metadatei.Sha().IsNull() {
		err = errors.Errorf("null sha")
		return
	}

	wg := iter.MakeErrorWaitGroupParallel()

	wg.Do(s.MakeFuncSaveOneVerweise(o))

	wg.Do(s.makeWriteMetadateiFunc(
		s.standort.DirVerzeichnisseMetadateiKennungMutter(),
		objekte_format.Formats.MetadateiKennungMutter(),
		o,
		o.Metadatei.Sha(),
	))

	wg.Do(s.makeWriteMetadateiFunc(
		s.standort.DirVerzeichnisseMetadatei(),
		objekte_format.Formats.Metadatei(),
		o,
		&o.Metadatei.SelbstMetadatei,
	))

	return wg.GetError()
}
