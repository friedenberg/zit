package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type Clean struct {
	force                     bool
	includeRecognizedAkten    bool
	includeRecognizedZettelen bool
	includeMutter             bool
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
				&c.includeMutter,
				"include-mutter",
				false,
				"remove Objekten in working directory if they match their Mutter",
			)

			f.BoolVar(
				&c.includeRecognizedAkten,
				"recognized-akten",
				false,
				"remove Akten in working directory or args that are recognized",
			)

			f.BoolVar(
				&c.includeRecognizedZettelen,
				"recognized-zettelen",
				false,
				"remove Zetteln in working directory or args that are recognized",
			)

			return c
		},
	)
}

func (c Clean) DefaultGattungen() ids.Genre {
	return ids.MakeGenre(genres.TrueGenre()...)
}

func (c Clean) shouldClean(
	u *umwelt.Umwelt,
	co sku.CheckedOutLike,
	eqwk *query.Group,
) bool {
	if c.force {
		return true
	}

	state := co.GetState()

	switch state {
	case checked_out_state.StateExistsAndSame:
		return true

	case checked_out_state.StateRecognized:
		return eqwk.IncludeRecognized
	}

	if c.includeMutter {
		mutter, err := u.GetStore().GetVerzeichnisse().ReadOneEnnui(
			co.GetSku().Metadatei.Mutter(),
		)

		errors.PanicIfError(err)

		if mutter != nil &&
			metadatei.EqualerSansTai.Equals(
				&co.GetSkuExternalLike().GetSku().Metadatei,
				&mutter.Metadatei,
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
	u *umwelt.Umwelt,
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

			// ui.Debug().Print(co)

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
