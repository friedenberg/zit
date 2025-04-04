package commands

import (
	"flag"
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_ptr"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register(
		"checkin-blob",
		&CheckinBlob{
			NewTags: collections_ptr.MakeFlagCommas[ids.Tag](
				collections_ptr.SetterPolicyAppend,
			),
		},
	)
}

type CheckinBlob struct {
	command_components.LocalWorkingCopy

	Delete  bool
	NewTags collections_ptr.Flag[ids.Tag, *ids.Tag]
}

func (cmd *CheckinBlob) SetFlagSet(f *flag.FlagSet) {
	f.BoolVar(&cmd.Delete, "delete", false, "the checked-out file")
	f.Var(
		cmd.NewTags,
		"new-tags",
		"comma-separated tags (will replace existing tags)",
	)
}

func (cmd CheckinBlob) Run(req command.Request) {
	args := req.PopArgs()

	localWorkingCopy := cmd.MakeLocalWorkingCopy(req)

	if len(args)%2 != 0 {
		req.CancelWithErrorf(
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
				req.CancelWithError(err)
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

			if zettels[i], err = localWorkingCopy.GetStore().ReadTransactedFromObjectId(
				p.ZettelId,
			); err != nil {
				req.CancelWithError(err)
			}
		}
	}

	for i, p := range pairs {
		var ow sha.WriteCloser

		{
			var err error

			if ow, err = localWorkingCopy.GetEnvRepo().BlobWriter(); err != nil {
				req.CancelWithError(err)
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
					req.CancelWithError(err)
				}
			}

			defer req.MustClose(f)

			if _, err := io.Copy(ow, f); err != nil {
				req.CancelWithError(err)
			}

			if err := ow.Close(); err != nil {
				req.CancelWithError(err)
			}

			{
				var err error

				if zettels[i], err = localWorkingCopy.GetStore().ReadTransactedFromObjectId(
					p.ZettelId,
				); err != nil {
					req.CancelWithError(err)
				}
			}

			zettels[i].SetBlobSha(ow.GetShaLike())

		case shaError == nil:
			zettels[i].SetBlobSha(&as)

		default:
			req.CancelWithError(errors.Errorf("argument is neither sha nor path"))
		}

		if cmd.NewTags.Len() > 0 {
			m := zettels[i].GetMetadata()
			m.SetTags(cmd.NewTags)
		}
	}

	req.Must(localWorkingCopy.Lock)

	for _, z := range zettels {
		if err := localWorkingCopy.GetStore().CreateOrUpdateDefaultProto(
			z,
			sku.StoreOptions{
				MergeCheckedOut: true,
			},
		); err != nil {
			req.CancelWithError(err)
		}
	}

	req.Must(localWorkingCopy.Unlock)
}
