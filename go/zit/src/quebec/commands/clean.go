package commands

import (
	"flag"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
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

func (c Clean) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung(gattung.TrueGattung()...)
}

func (c Clean) shouldClean(u *umwelt.Umwelt, co *sku.CheckedOut) bool {
	ui.Log().Print(co)
	if co.State == checked_out_state.StateExistsAndSame {
		return true
	}

	if c.force {
		return true
	}

	if c.includeMutter {
		mutter, err := u.GetStore().GetVerzeichnisse().ReadOneEnnui(
			co.Internal.Metadatei.Mutter(),
		)

		errors.PanicIfError(err)

		if mutter != nil &&
			metadatei.EqualerSansTai.Equals(&co.External.Metadatei, &mutter.Metadatei) {
			return true
		}
	}

	return false
}

func (c Clean) ModifyBuilder(
	b *query.Builder,
) {
	b.WithHidden(nil)
}

func (c Clean) RunWithQuery(
	u *umwelt.Umwelt,
	qg *query.Group,
) (err error) {
	fds := fd.MakeMutableSet()
	l := &sync.Mutex{}

	for _, d := range u.GetStore().GetCwdFiles().EmptyDirectories {
		fds.Add(d)
	}

	if err = u.GetStore().ReadFiles(
		qg,
		func(co *sku.CheckedOut) (err error) {
			if !c.shouldClean(u, co) {
				return
			}

			e := co.External

			l.Lock()
			defer l.Unlock()

			fds.Add(e.GetObjekteFD())
			fds.Add(e.GetAkteFD())

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.markUnsureAktenForRemovalIfNecessary(u, qg, fds.Add); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.markUnsureZettelenForRemovalIfNecessary(
		u,
		qg,
		fds.Add,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	deleteOp := user_ops.DeleteCheckout{
		Umwelt: u,
	}

	if err = deleteOp.Run(fds); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Clean) markUnsureAktenForRemovalIfNecessary(
	u *umwelt.Umwelt,
	qg *query.Group,
	add schnittstellen.FuncIter[*fd.FD],
) (err error) {
	if !c.includeRecognizedAkten {
		return
	}

	if err = qg.GetExplicitCwdFDs().Each(
		u.GetStore().GetCwdFiles().MarkUnsureAkten,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := u.PrinterCheckedOut()
	var l sync.Mutex

	if err = u.GetStore().ReadAllMatchingAkten(
		qg,
		u.GetStore().GetCwdFiles().UnsureAkten,
		func(fd *fd.FD, z *sku.Transacted) (err error) {
			if z == nil {
				err = u.PrinterFileNotRecognized()(fd)
				return
			}

			os := sha.Make(z.GetObjekteSha())
			as := sha.Make(z.GetAkteSha())

			fr := sku.GetCheckedOutPool().Get()
			defer sku.GetCheckedOutPool().Put(fr)

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

func (c Clean) markUnsureZettelenForRemovalIfNecessary(
	u *umwelt.Umwelt,
	qg *query.Group,
	add schnittstellen.FuncIter[*fd.FD],
) (err error) {
	if !c.includeRecognizedZettelen {
		return
	}

	p := u.PrinterCheckedOut()
	var l sync.Mutex

	if err = u.GetStore().QueryUnsure(
		qg,
		store.UnsureMatchOptions{
			UnsureMatchType: store.UnsureMatchTypeMetadateiSansTaiHistory,
		},
		func(
			mt store.UnsureMatchType,
			sk *sku.Transacted,
			existing sku.CheckedOutMutableSet,
		) (err error) {
			if err = existing.Each(
				func(fr *sku.CheckedOut) (err error) {
					fr.State = checked_out_state.StateRecognized

					if err = p(fr); err != nil {
						err = errors.Wrap(err)
						return
					}

					l.Lock()
					defer l.Unlock()

					if !fr.External.FDs.Objekte.IsEmpty() {
						if err = add(&fr.External.FDs.Objekte); err != nil {
							err = errors.Wrap(err)
							return
						}
					}

					if !fr.External.FDs.Akte.IsEmpty() {
						if err = add(&fr.External.FDs.Akte); err != nil {
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
