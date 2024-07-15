package inventory_list

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/golf/sha_probe_index"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type probe_store struct {
	fs_home                   fs_home.Home
	persistentMetadateiFormat object_inventory_format.Format
	object_id_probe           sha_probe_index.Ennui
	options                   object_inventory_format.Options
}

func (s *probe_store) Initialize(
	fs_home fs_home.Home,
	persistentMetadateiFormat object_inventory_format.Format,
	options object_inventory_format.Options,
) (err error) {
	s.fs_home = fs_home
	s.persistentMetadateiFormat = persistentMetadateiFormat
	s.options = options

	if s.object_id_probe, err = sha_probe_index.MakeNoDuplicates(
		s.fs_home,
		s.fs_home.DirVerzeichnisseVerweise(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *probe_store) ReadOneEnnui(sh *sha.Sha) (sk *sku.Transacted, err error) {
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

	if err = sk.CalculateObjectShas(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *probe_store) ReadOneObjectIdSha(k ids.IdLike) (sh *sha.Sha, err error) {
	left := sha.FromString(k.String())
	defer sha.GetPool().Put(left)

	if sh, err = s.object_id_probe.ReadOne(left); err != nil {
		err = errors.Wrapf(err, "object id: %q, Left: %s", k, left)
		return
	}

	return
}

func (s *probe_store) ReadOneObjectId(k ids.IdLike) (sk *sku.Transacted, err error) {
	sh, err := s.ReadOneObjectIdSha(k)
	defer sha.GetPool().Put(sh)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if sk, err = s.ReadOneEnnui(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !sh.Equals(sk.Metadata.Sha()) {
		err = errors.Errorf(
			"expected sha %q but got %q",
			sh,
			sk.Metadata.Sha(),
		)

		return
	}

	return
}

func (s *probe_store) makeWriteMetadateiFunc(
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

func (s *probe_store) MakeFuncSaveOneVerweise(o *sku.Transacted) func() error {
	return func() (err error) {
		k := o.GetObjectId()
		sh := sha.FromString(k.String())
		defer sha.GetPool().Put(sh)

		if err = s.object_id_probe.AddSha(sh, o.Metadata.Sha()); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (s *probe_store) WriteOneObjekteMetadatei(o *sku.Transacted) (err error) {
	if o.Metadata.Sha().IsNull() {
		err = errors.Errorf("null sha")
		return
	}

	wg := iter.MakeErrorWaitGroupParallel()

	wg.Do(s.MakeFuncSaveOneVerweise(o))

	wg.Do(s.makeWriteMetadateiFunc(
		s.fs_home.DirVerzeichnisseMetadateiKennungMutter(),
		object_inventory_format.Formats.MetadateiKennungMutter(),
		o,
		o.Metadata.Sha(),
	))

	wg.Do(s.makeWriteMetadateiFunc(
		s.fs_home.DirVerzeichnisseMetadatei(),
		object_inventory_format.Formats.Metadatei(),
		o,
		&o.Metadata.SelfMetadata,
	))

	return wg.GetError()
}
