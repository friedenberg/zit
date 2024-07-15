package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Status struct{}

func init() {
	registerCommandWithQuery(
		"status",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Status{}

			return c
		},
	)
}

func (c Status) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.TrueGenre()...)
}

func (c Status) ModifyBuilder(
	b *query.Builder,
) {
	b.WithHidden(nil)
}

func (c Status) RunWithQuery(
	u *env.Env,
	eqwk *query.Group,
) (err error) {
	pcol := u.PrinterCheckedOutForKasten(eqwk.RepoId)

	if err = u.GetStore().QueryCheckedOut(
		eqwk,
		func(co sku.CheckedOutLike) (err error) {
			if err = pcol(co); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.GetStore().QueryUnsure(
		eqwk,
		sku.UnsureMatchOptions{
			UnsureMatchType: sku.UnsureMatchTypeMetadataWithoutTaiHistory | sku.UnsureMatchTypeDescription,
		},
		u.PrinterMatching(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO [ces/mew] switch to kasten parsing ID's before body
	// if err = qg.GetExplicitCwdFDs().Each(
	// 	u.GetStore().GetCwdFiles().MarkUnsureAkten,
	// ); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	p := u.PrinterCheckedOutForKasten(eqwk.RepoId)
	qg := eqwk

	// TODO [cot/gl !task project-2021-zit-kasten today zz-inbox] move unsure akten and untracked into kasten interface and store_fs
	if err = u.GetStore().QueryAllMatchingAkten(
		qg,
		u.GetStore().GetCwdFiles().GetUnsureAkten(),
		func(fd *fd.FD, z *sku.Transacted) (err error) {
			if z == nil {
				err = u.PrinterFileNotRecognized()(fd)
				return
			}

			os := sha.Make(z.GetObjectSha())
			as := sha.Make(z.GetBlobSha())

			fr := store_fs.GetCheckedOutPool().Get()
			defer store_fs.GetCheckedOutPool().Put(fr)

			fr.State = checked_out_state.StateRecognized

			if err = fr.Internal.SetFromSkuLike(z); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = fr.External.SetFromSkuLike(z); err != nil {
				err = errors.Wrap(err)
				return
			}

			fr.External.FDs.Blob.ResetWith(fd)
			fr.External.SetBlobSha(as)

			if err = fr.External.SetObjectSha(os); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = p(fr); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
