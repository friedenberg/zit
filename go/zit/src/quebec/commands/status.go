package commands

import (
	"flag"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
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

func (c Status) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung(gattung.TrueGattung()...)
}

func (c Status) ModifyBuilder(
	b *query.Builder,
) {
	b.WithHidden(nil)
}

func (c Status) RunWithQuery(
	u *umwelt.Umwelt,
	qg *query.Group,
) (err error) {
	pcol := u.PrinterCheckedOut()

	if err = u.GetStore().ReadFiles(
		qg,
		func(co *sku.CheckedOut) (err error) {
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

	selbstMetadateiSansTaiToZettels := make(
		map[sha.Bytes]sku.CheckedOutMutableSet,
		u.GetStore().GetCwdFiles().Len(),
	)

	bezToZettels := make(
		map[string]sku.CheckedOutMutableSet,
		u.GetStore().GetCwdFiles().Len(),
	)

	var l sync.Mutex

	if err = u.GetStore().ReadFilesUnsure(
		qg,
		func(co *sku.CheckedOut) (err error) {
			if err = pcol(co); err != nil {
				err = errors.Wrap(err)
				return
			}

			sh := &co.External.Metadatei.Shas.SelbstMetadateiSansTai

			if sh.IsNull() {
				return
			}

			l.Lock()
			defer l.Unlock()

			clone := sku.GetCheckedOutPool().Get()
			sku.CheckedOutResetter.ResetWith(clone, co)

			{
				k := sh.GetBytes()
				existing, ok := selbstMetadateiSansTaiToZettels[k]

				if !ok {
					existing = sku.MakeCheckedOutMutableSet()
				}

				if err = existing.Add(clone); err != nil {
					err = errors.Wrap(err)
					return
				}

				selbstMetadateiSansTaiToZettels[k] = existing
			}

			{
				k := co.External.Metadatei.Bezeichnung.String()
				existing, ok := bezToZettels[k]

				if !ok {
					existing = sku.MakeCheckedOutMutableSet()
				}

				if err = existing.Add(clone); err != nil {
					err = errors.Wrap(err)
					return
				}

				bezToZettels[k] = existing
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := u.PrinterCheckedOut()

	qg.SetIncludeHistory()

	if len(selbstMetadateiSansTaiToZettels) > 0 || len(bezToZettels) > 0 {
		if err = u.GetStore().QueryWithoutCwd(
			qg,
			func(sk *sku.Transacted) (err error) {
				sh := &sk.Metadatei.Shas.SelbstMetadateiSansTai

				if sh.IsNull() {
					return
				}

				printMatching := func(existing sku.CheckedOutMutableSet) (err error) {
					if err = existing.Each(
						func(co *sku.CheckedOut) (err error) {
							co.State = checked_out_state.StateRecognized

							if err = co.Internal.SetFromSkuLike(sk); err != nil {
								err = errors.Wrap(err)
								return
							}

							if err = p(co); err != nil {
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

				{
					k := sh.GetBytes()
					existing, ok := selbstMetadateiSansTaiToZettels[k]

					if !ok {
						return
					}

					if err = printMatching(existing); err != nil {
						err = errors.Wrap(err)
						return
					}
				}

				{
					k := sk.Metadatei.Bezeichnung.String()
					existing, ok := bezToZettels[k]

					if !ok {
						return
					}

					if err = printMatching(existing); err != nil {
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
	}

	if err = qg.GetExplicitCwdFDs().Each(
		u.GetStore().GetCwdFiles().MarkUnsureAkten,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

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
			fr.External.SetAkteSha(as)

			if err = fr.External.SetObjekteSha(os); err != nil {
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
