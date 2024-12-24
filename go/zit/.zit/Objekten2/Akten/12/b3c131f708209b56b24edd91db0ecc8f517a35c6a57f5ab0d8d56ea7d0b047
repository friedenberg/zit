package commands

import (
	"flag"
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_ptr"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type CheckinBlob struct {
	Delete  bool
	NewTags collections_ptr.Flag[ids.Tag, *ids.Tag]
}

func init() {
	registerCommand(
		"checkin-blob",
		func(f *flag.FlagSet) Command {
			c := &CheckinBlob{
				NewTags: collections_ptr.MakeFlagCommas[ids.Tag](
					collections_ptr.SetterPolicyAppend,
				),
			}

			f.BoolVar(&c.Delete, "delete", false, "the checked-out file")
			f.Var(
				c.NewTags,
				"new-tags",
				"comma-separated tags (will replace existing tags)",
			)

			return c
		},
	)
}

func (c CheckinBlob) Run(u *env.Local, args ...string) (err error) {
	if len(args)%2 != 0 {
		err = errors.Errorf(
			"arguments must come in pairs of zettel id and blob path",
		)
		return
	}

	type externalBlobPair struct {
		*ids.ZettelId
		path string
	}

	pairs := make([]externalBlobPair, len(args)/2)

	// transform args into pairs of hinweis and filepaths
	for i, p := range pairs {
		hs := args[i*2]
		ap := args[(i*2)+1]

		if p.ZettelId, err = ids.MakeZettelId(hs); err != nil {
			err = errors.Wrap(err)
			return
		}

		p.path = ap
		pairs[i] = p
	}

	zettels := make([]*sku.Transacted, len(pairs))

	// iterate through pairs and read current zettel
	for i, p := range pairs {
		if zettels[i], err = u.GetStore().ReadTransactedFromObjectId(
			p.ZettelId,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for i, p := range pairs {
		var ow sha.WriteCloser

		if ow, err = u.GetDirectoryLayout().BlobWriter(); err != nil {
			err = errors.Wrap(err)
			return
		}

		var as sha.Sha

		shaError := as.Set(p.path)

		switch {
		case files.Exists(p.path):
			var f *os.File

			if f, err = files.Open(p.path); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.DeferredCloser(&err, f)

			if _, err = io.Copy(ow, f); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = ow.Close(); err != nil {
				err = errors.Wrap(err)
				return
			}

			if zettels[i], err = u.GetStore().ReadTransactedFromObjectId(
				p.ZettelId,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			zettels[i].SetBlobSha(ow.GetShaLike())

		case shaError == nil:
			zettels[i].SetBlobSha(&as)

		default:
			err = errors.Errorf("argument is neither sha nor path")
			return
		}

		if c.NewTags.Len() > 0 {
			m := zettels[i].GetMetadata()
			m.SetTags(c.NewTags)
		}
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	for _, z := range zettels {
		if err = u.GetStore().CreateOrUpdate(
			z,
			object_mode.Make(
				object_mode.ModeMergeCheckedOut,
			),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
