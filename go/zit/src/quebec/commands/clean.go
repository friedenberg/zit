package commands

import (
	"flag"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type Clean struct {
	force                     bool
	includeRecognizedAkten    bool
	includeRecognizedZettelen bool
	includeMutter             bool
}

func init() {
	registerCommandWithExternalQuery(
		"clean",
		func(f *flag.FlagSet) CommandWithExternalQuery {
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

func (c Clean) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung(gattung.TrueGattung()...)
}

func (c Clean) shouldClean(
	u *umwelt.Umwelt,
	co sku.CheckedOutLike,
	eqwk sku.ExternalQuery,
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

func (c Clean) RunWithExternalQuery(
	u *umwelt.Umwelt,
	eqwk sku.ExternalQuery,
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

// TODO [radi/kof !task project-2021-zit-features zz-inbox] add support for kasten in checkouts and external
func (c Clean) markUnsureAktenForRemovalIfNecessary(
	u *umwelt.Umwelt,
	qg *query.Group,
	k kennung.Kasten,
	add schnittstellen.FuncIter[*fd.FD],
) (err error) {
	if !c.includeRecognizedAkten {
		return
	}

	// TODO [ces/mew] switch to kasten parsing ID's before body
	// if err = qg.GetExplicitCwdFDs().Each(
	// 	u.GetStore().GetCwdFiles().MarkUnsureAkten,
	// ); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	p := u.PrinterCheckedOutForKasten(k)
	var l sync.Mutex

	// TODO create a new query group for all of history
	qg.SetIncludeHistory()

	if err = u.GetStore().QueryAllMatchingAkten(
		qg,
		u.GetStore().GetCwdFiles().GetUnsureAkten(),
		func(fd *fd.FD, z *sku.Transacted) (err error) {
			if z == nil {
				err = u.PrinterFileNotRecognized()(fd)
				return
			}

			os := sha.Make(z.GetObjekteSha())
			as := sha.Make(z.GetAkteSha())

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

			fr.External.FDs.Akte.ResetWith(fd)

			if err = fr.External.SetAkteSha(as); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = fr.External.SetObjekteSha(os); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = p(fr); err != nil {
				err = errors.Wrap(err)
				return
			}

			l.Lock()
			defer l.Unlock()

			if err = add(fd); err != nil {
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

// TODO [radi/kof !task project-2021-zit-features zz-inbox] add support for kasten in checkouts and external
func (c Clean) markUnsureZettelenForRemovalIfNecessary(
	u *umwelt.Umwelt,
	qg *query.Group,
	k kennung.Kasten,
	add schnittstellen.FuncIter[*fd.FD],
) (err error) {
	if !c.includeRecognizedZettelen {
		return
	}

	p := u.PrinterCheckedOutForKasten(k)
	var l sync.Mutex

	if err = u.GetStore().QueryUnsure(
		qg,
		store.UnsureMatchOptions{
			UnsureMatchType: store.UnsureMatchTypeMetadateiSansTaiHistory,
		},
		func(
			mt store.UnsureMatchType,
			sk *sku.Transacted,
			existing sku.CheckedOutLikeMutableSet,
		) (err error) {
			if err = existing.Each(
				func(fr sku.CheckedOutLike) (err error) {
					if err = fr.SetState(
						checked_out_state.StateRecognized,
					); err != nil {
						err = errors.Wrap(err)
						return
					}

					if err = p(fr); err != nil {
						err = errors.Wrap(err)
						return
					}

					l.Lock()
					defer l.Unlock()

					// TODO add support for other checked out types
					cofs, ok := fr.(*store_fs.CheckedOut)

					if !ok {
						return
					}

					if !cofs.External.FDs.Objekte.IsEmpty() {
						if err = add(&cofs.External.FDs.Objekte); err != nil {
							err = errors.Wrap(err)
							return
						}
					}

					if !cofs.External.FDs.Akte.IsEmpty() {
						if err = add(&cofs.External.FDs.Akte); err != nil {
							err = errors.Wrap(err)
							return
						}
					}

					return
				},
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

	return
}
