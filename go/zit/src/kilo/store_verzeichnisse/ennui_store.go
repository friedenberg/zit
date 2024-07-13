package store_verzeichnisse

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/golf/sha_probe_index"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type ennuiStore struct {
	fs_home                   fs_home.Standort
	persistentMetadateiFormat object_inventory_format.Format
	ids                       sha_probe_index.Ennui
	options                   object_inventory_format.Options
}

func (s *ennuiStore) Initialize(
	fs_home fs_home.Standort,
	persistentMetadateiFormat object_inventory_format.Format,
	options object_inventory_format.Options,
) (err error) {
	s.fs_home = fs_home
	s.persistentMetadateiFormat = persistentMetadateiFormat
	s.options = options

	if s.ids, err = sha_probe_index.MakeNoDuplicates(
		s.fs_home,
		s.fs_home.DirVerzeichnisseVerweise(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *ennuiStore) ReadOneEnnui(sh *sha.Sha) (sk *sku.Transacted, err error) {
	var r sha.ReadCloser

	if r, err = s.fs_home.BlobReaderFrom(
		sh,
		s.fs_home.DirVerzeichnisseMetadateiKennungMutter(),
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

func (s *ennuiStore) ReadOneKennungSha(
	k interfaces.StringerGenreGetter,
) (sh *sha.Sha, err error) {
	left := sha.FromString(k.String())
	defer sha.GetPool().Put(left)

	if sh, err = s.ids.ReadOne(left); err != nil {
		err = errors.Wrapf(err, "Kennung: %q, Left: %s", k, left)
		return
	}

	return
}

func (s *ennuiStore) ReadOneKennung(
	k interfaces.StringerGenreGetter,
) (sk *sku.Transacted, err error) {
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
	fo object_inventory_format.FormatGeneric,
	o *sku.Transacted,
	expected *sha.Sha,
) interfaces.FuncError {
	return func() (err error) {
		var sw sha.WriteCloser

		if sw, err = s.fs_home.BlobWriterToLight(dir); err != nil {
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

		if err = s.ids.AddSha(sh, o.Metadatei.Sha()); err != nil {
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
		s.fs_home.DirVerzeichnisseMetadateiKennungMutter(),
		object_inventory_format.Formats.MetadateiKennungMutter(),
		o,
		o.Metadatei.Sha(),
	))

	wg.Do(s.makeWriteMetadateiFunc(
		s.fs_home.DirVerzeichnisseMetadatei(),
		object_inventory_format.Formats.Metadatei(),
		o,
		&o.Metadatei.SelbstMetadatei,
	))

	return wg.GetError()
}
