package commands

import (
	"flag"
	"sync"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/matcher_proto"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
	"code.linenisgreat.com/zit/src/papa/user_ops"
)

type Clean struct {
	force             bool
	includeRecognized bool
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
				&c.includeRecognized,
				"recognized",
				false,
				"remove Akten in working directory or args that are recognized",
			)

			return c
		},
	)
}

func (c Clean) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung(gattung.TrueGattung()...)
}

func (c Clean) RunWithQuery(
	u *umwelt.Umwelt,
	ms matcher_proto.QueryGroup,
) (err error) {
	fds := fd.MakeMutableSet()
	l := &sync.Mutex{}

	for _, d := range u.StoreUtil().GetCwdFiles().EmptyDirectories {
		fds.Add(d)
	}

	if err = u.StoreObjekten().ReadFiles(
		matcher_proto.MakeFuncReaderTransactedLikePtr(ms, u.StoreObjekten().QueryWithoutCwd),
		iter.MakeChain(
			matcher_proto.MakeFilterFromQuery(ms),
			func(co *sku.CheckedOut) (err error) {
				if co.State != checked_out_state.StateExistsAndSame && !c.force {
					return
				}

				e := co.External

				l.Lock()
				defer l.Unlock()

				fds.Add(e.GetObjekteFD())
				fds.Add(e.GetAkteFD())

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.markUnsureAktenForRemovalIfNecessary(u, ms, fds.Add); err != nil {
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
	q matcher_proto.QueryGroup,
	add schnittstellen.FuncIter[*fd.FD],
) (err error) {
	if !c.includeRecognized {
		return
	}

	if err = q.GetExplicitCwdFDs().Each(
		u.StoreUtil().GetCwdFiles().MarkUnsureAkten,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := u.PrinterCheckedOut()
	var l sync.Mutex

	if err = u.StoreObjekten().ReadAllMatchingAkten(
		u.StoreUtil().GetCwdFiles().UnsureAkten,
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
