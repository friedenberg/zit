package user_ops

import (
	"io"
	"sort"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type ZettelFromExternalBlob struct {
	*env.Env
	sku.Proto
	// TODO switch to using ObjekteOptions
	Filter script_value.ScriptValue
	Delete bool
	Dedupe bool
}

func (c ZettelFromExternalBlob) Run(
	fdSet fd.Set,
) (results sku.TransactedMutableSet, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, c.Unlock)

	results = sku.MakeTransactedMutableSet()
	toDelete := fd.MakeMutableSet()

	fds := make(map[sha.Bytes][]*fd.FD)

	for _, fd := range iter.SortedValues(fdSet) {
		if err = c.addToMapAndWriteToBlobStore(fd, fds); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	toCreate := make([]*store_fs.External, 0, len(fds))

	for _, fdsForSha := range fds {
		sort.Slice(fdsForSha, func(i, j int) bool {
			return fdsForSha[i].String() < fdsForSha[j].String()
		})

		var z *store_fs.External

		if z, err = c.createZettelForBlobs(fdsForSha); err != nil {
			err = errors.Wrap(err)
			return
		}

		toCreate = append(toCreate, z)

		if !c.Delete {
			continue
		}

		for _, fd := range fdsForSha {
			if err = toDelete.Add(fd); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	sort.Slice(
		toCreate,
		func(i, j int) bool {
			return toCreate[i].GetBlobFD().String() < toCreate[j].GetBlobFD().String()
		},
	)

	for _, z := range toCreate {
		if z.Metadata.IsEmpty() {
			return
		}

		if err = c.GetStore().CreateOrUpdate(
			z,
			objekte_mode.ModeApplyProto,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		// TODO switch to using ObjekteOptions
		if c.Proto.Apply(z, genres.Zettel) {
			if err = c.GetStore().CreateOrUpdate(
				z.GetSku(),
				objekte_mode.ModeEmpty,
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		results.Add(z.GetSku())
	}

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO move to umwelt
	dp := c.Env.PrinterFDDeleted()

	err = toDelete.Each(
		func(f *fd.FD) (err error) {
			if err = c.Env.GetFSHome().Delete(f.GetPath()); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = dp(f); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *ZettelFromExternalBlob) addToMapAndWriteToBlobStore(
	f *fd.FD,
	fds map[sha.Bytes][]*fd.FD,
) (err error) {
	var r io.Reader

	if r, err = c.Filter.Run(f.GetPath()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, &c.Filter)

	var blobWriter sha.WriteCloser

	if blobWriter, err = c.GetFSHome().BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, blobWriter)

	if _, err = io.Copy(blobWriter, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	f.SetShaLike(blobWriter.GetShaLike())

	key := sha.Make(f.GetShaLike()).GetBytes()
	existing := fds[key]
	existing = append(existing, f)
	fds[key] = existing

	return
}

func (c *ZettelFromExternalBlob) createZettelForBlobs(
	blobFDs []*fd.FD,
) (z *store_fs.External, err error) {
	// TODO handle other FD's
	blobFD := blobFDs[0]
	z = store_fs.GetExternalPool().Get()

	z.FDs.Blob.ResetWith(blobFD)

	z.Transacted.ObjectId.SetGenre(genres.Zettel)

	if err = c.Proto.ApplyWithBlobFD(z, blobFD); err != nil {
		err = errors.Wrap(err)
		return
	}

	z.SetBlobSha(blobFD.GetShaLike())

	return
}
