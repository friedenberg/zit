package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type CommitOrganizeFile struct {
	*umwelt.Umwelt
	OutputJSON bool
}

type CommitOrganizeFileResults = organize_text.Changes

func (c CommitOrganizeFile) ApplyToText(
	u *umwelt.Umwelt,
	t *organize_text.Text,
) (err error) {
	if u.GetKonfig().PrintOptions.PrintEtikettenAlways {
		return
	}

	if err = t.Transacted.Each(
		func(sk *sku.Transacted) (err error) {
			if sk.Metadatei.Description.IsEmpty() {
				return
			}

			sk.Metadatei.ResetEtiketten()

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op CommitOrganizeFile) Run(
	u *umwelt.Umwelt,
	a, b *organize_text.Text,
	original sku.TransactedSet,
) (results CommitOrganizeFileResults, err error) {
	if results, err = op.run(u, a, b, original); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op CommitOrganizeFile) run(
	u *umwelt.Umwelt,
	a, b *organize_text.Text,
	original sku.TransactedSet,
) (cs CommitOrganizeFileResults, err error) {
	if err = op.ApplyToText(u, a); err != nil {
		err = errors.Wrap(err)
		return
	}

	if cs, err = organize_text.ChangesFrom(
		a,
		b,
		original,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO move to bestandsaufnahme flush
	// if cs.GetAddedUnnamed().Len() == 0 && cs.GetAddedNamed().Len() == 0 {
	// 	errors.Err().Print("no changes")
	// 	return
	// }

	if err = cs.Changed.Each(
		func(changed *sku.Transacted) (err error) {
			if err = u.GetStore().CreateOrUpdateFromTransacted(
				changed,
				objekte_mode.Make(
					objekte_mode.ModeMergeCheckedOut,
				),
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
