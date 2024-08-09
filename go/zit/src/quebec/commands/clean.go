package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Clean struct {
	force                    bool
	includeRecognizedBlobs   bool
	includeRecognizedZettels bool
	includeParent            bool
}

func init() {
	registerCommandWithQuery(
		"clean",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Clean{}

			f.BoolVar(
				&c.force,
				"force",
				false,
				"remove Objekten in working directory even if they have changes",
			)

			f.BoolVar(
				&c.includeParent,
				"include-mutter",
				false,
				"remove Objekten in working directory if they match their Mutter",
			)

			f.BoolVar(
				&c.includeRecognizedBlobs,
				"recognized-blobs",
				false,
				"remove blobs in working directory or args that are recognized",
			)

			f.BoolVar(
				&c.includeRecognizedZettels,
				"recognized-zettelen",
				false,
				"remove Zetteln in working directory or args that are recognized",
			)

			return c
		},
	)
}

func (c Clean) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.TrueGenre()...)
}

func (c Clean) shouldClean(
	u *env.Env,
	co sku.CheckedOutLike,
	eqwk *query.Group,
) bool {
	if c.force {
		return true
	}

	state := co.GetState()

	switch state {
	case checked_out_state.ExistsAndSame:
		return true

	case checked_out_state.Recognized:
		return !eqwk.ExcludeRecognized
	}

	if c.includeParent {
		mutter := sku.GetTransactedPool().Get()
		defer sku.GetTransactedPool().Put(mutter)

		err := u.GetStore().GetStreamIndex().ReadOneObjectId(
			co.GetSku().GetObjectId().String(),
			mutter,
		)

		errors.PanicIfError(err)

		if object_metadata.EqualerSansTai.Equals(
			&co.GetSkuExternalLike().GetSku().Metadata,
			&mutter.Metadata,
		) {
			return true
		}
	}

	return false
}

func (c Clean) ModifyBuilder(b *query.Builder) {
	b.WithHidden(nil)
}

func (c Clean) RunWithQuery(
	u *env.Env,
	eqwk *query.Group,
) (err error) {
	fds := fd.MakeMutableSet()

	u.Lock()
	defer errors.Deferred(&err, u.Unlock)

	// TODO [radi/kof !task project-2021-zit-features zz-inbox] add support for kasten in checkouts and external
	if err = u.GetStore().GetCwdFiles().GetEmptyDirectories().Each(
		fds.Add,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.GetStore().QueryCheckedOut(
		eqwk,
		func(co sku.CheckedOutLike) (err error) {
			if !c.shouldClean(u, co, eqwk) {
				return
			}

			if err = u.GetStore().DeleteCheckout(co); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO [radi/kof !task project-2021-zit-features zz-inbox] add support for kasten in checkouts and external
	// if err = c.markUnsureAktenForRemovalIfNecessary(u, eqwk, fds.Add); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	// TODO [radi/kof !task project-2021-zit-features zz-inbox] add support for kasten in checkouts and external
	// if err = c.markUnsureZettelenForRemovalIfNecessary(
	// 	u,
	// 	eqwk,
	// 	fds.Add,
	// ); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	if err = u.DeleteFiles(fds); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
