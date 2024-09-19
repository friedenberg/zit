package user_ops

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type ZettelFromExternalBlob struct {
	*env.Env
	sku.Proto
	// TODO switch to using ObjekteOptions
	Filter     script_value.ScriptValue
	Delete     bool
	AllowDupes bool
}

func (c ZettelFromExternalBlob) Run(
	qg *query.Group,
) (results sku.TransactedMutableSet, err error) {
	results = sku.MakeTransactedMutableSet()
	toDelete := fd.MakeMutableSet()

	if err = c.GetStore().QueryCheckedOut(
		qg,
		func(col sku.CheckedOutLike) (err error) {
			// TODO support other repos
			cofs := col.(*store_fs.CheckedOut)
			z := col.GetSkuExternalLike().GetSku()

			if z.Metadata.IsEmpty() {
				return
			}

			if err = c.GetStore().GetCwdFiles().UpdateDescriptionFromBlobs(&cofs.External); err != nil {
				err = errors.Wrap(err)
				return
			}

			z.ObjectId.Reset()

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

			if err = results.Add(z.GetSku()); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = c.Env.GetStore().DeleteExternalLike(
				ids.RepoId{},
				&cofs.External,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
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
) (z *sku.External, err error) {
	// TODO handle other FD's
	blobFD := blobFDs[0]
	z = store_fs.GetExternalPool().Get()

	if err = c.GetStore().GetCwdFiles().SetBlobOrError(z, blobFD); err != nil {
		err = errors.Wrap(err)
		return
	}

	z.Transacted.ObjectId.SetGenre(genres.Zettel)

	if err = c.Proto.ApplyWithBlobFD(&z.Transacted, blobFD); err != nil {
		err = errors.Wrap(err)
		return
	}

	z.SetBlobSha(blobFD.GetShaLike())

	return
}
