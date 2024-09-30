package inventory_list

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/golf/sha_probe_index"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type probe_index struct {
	fs_home                  fs_home.Home
	persistentMetadataFormat object_inventory_format.Format
	object_id_probe          sha_probe_index.Index
	options                  object_inventory_format.Options
}

func (s *probe_index) Initialize(
	fs_home fs_home.Home,
	persistentMetadataFormat object_inventory_format.Format,
	options object_inventory_format.Options,
) (err error) {
	s.fs_home = fs_home
	s.persistentMetadataFormat = persistentMetadataFormat
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

// TODO replace this with a much better and more efficient object lookup store
// (for writes, currently this makes one file per object which is extremely
// inefficient)
func (s *probe_index) makeWriteMetadataFunc(
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

		if _, err = fo.WriteMetadataTo(sw, o); err != nil {
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

func (s *probe_index) MakeFuncSaveOneObjectId(o *sku.Transacted) func() error {
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

func (s *probe_index) WriteOneObject(o *sku.Transacted) (err error) {
	if o.Metadata.Sha().IsNull() {
		err = errors.Errorf("null sha")
		return
	}

	wg := quiter.MakeErrorWaitGroupParallel()

	wg.Do(s.MakeFuncSaveOneObjectId(o))

	wg.Do(s.makeWriteMetadataFunc(
		s.fs_home.DirVerzeichnisseMetadataObjectIdParent(),
		object_inventory_format.Formats.MetadataObjectIdParent(),
		o,
		o.Metadata.Sha(),
	))

	wg.Do(s.makeWriteMetadataFunc(
		s.fs_home.DirVerzeichnisseMetadata(),
		object_inventory_format.Formats.Metadata(),
		o,
		&o.Metadata.SelfMetadata,
	))

	return wg.GetError()
}
