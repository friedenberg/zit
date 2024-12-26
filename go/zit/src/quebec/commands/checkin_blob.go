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
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type CheckinBlob struct {
	Delete  bool
	NewTags collections_ptr.Flag[ids.Tag, *ids.Tag]
}

func init() {
	registerCommand(
		"checkin-blob",
		func(f *flag.FlagSet) CommandWithRepo {
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

func (c CheckinBlob) RunWithRepo(u *repo_local.Repo, args ...string) {
	if len(args)%2 != 0 {
		u.CancelWithErrorf(
			"arguments must come in pairs of zettel id and blob path",
		)
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

		{
			var err error

			if p.ZettelId, err = ids.MakeZettelId(hs); err != nil {
				u.CancelWithError(err)
			}
		}

		p.path = ap
		pairs[i] = p
	}

	zettels := make([]*sku.Transacted, len(pairs))

	// iterate through pairs and read current zettel
	for i, p := range pairs {
		{
			var err error

			if zettels[i], err = u.GetStore().ReadTransactedFromObjectId(
				p.ZettelId,
			); err != nil {
				u.CancelWithError(err)
			}
		}
	}

	for i, p := range pairs {
		var ow sha.WriteCloser

		{
			var err error

			if ow, err = u.GetRepoLayout().BlobWriter(); err != nil {
				u.CancelWithError(err)
			}
		}

		var as sha.Sha

		shaError := as.Set(p.path)

		switch {
		case files.Exists(p.path):
			var f *os.File

			{
				var err error

				if f, err = files.Open(p.path); err != nil {
					u.CancelWithError(err)
				}
			}

			defer u.MustClose(f)

			if _, err := io.Copy(ow, f); err != nil {
				u.CancelWithError(err)
			}

			if err := ow.Close(); err != nil {
				u.CancelWithError(err)
			}

			{
				var err error

				if zettels[i], err = u.GetStore().ReadTransactedFromObjectId(
					p.ZettelId,
				); err != nil {
					u.CancelWithError(err)
				}
			}

			zettels[i].SetBlobSha(ow.GetShaLike())

		case shaError == nil:
			zettels[i].SetBlobSha(&as)

		default:
			u.CancelWithError(errors.Errorf("argument is neither sha nor path"))
		}

		if c.NewTags.Len() > 0 {
			m := zettels[i].GetMetadata()
			m.SetTags(c.NewTags)
		}
	}

	if err := u.Lock(); err != nil {
		u.CancelWithError(err)
	}

	defer u.Must(u.Unlock)

	for _, z := range zettels {
		if err := u.GetStore().CreateOrUpdate(
			z,
			object_mode.Make(
				object_mode.ModeMergeCheckedOut,
			),
		); err != nil {
			u.CancelWithError(err)
		}
	}
}
