package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type CommitOrganizeFile struct {
	*env.Env
	OutputJSON bool
}

type CommitOrganizeFileResults = organize_text.Changes

func (c CommitOrganizeFile) ApplyToText(
	u *env.Env,
	t *organize_text.Text,
) (err error) {
	if u.GetConfig().PrintOptions.PrintTagsAlways {
		return
	}

	if err = t.Transacted.Each(
		func(el sku.ExternalLike) (err error) {
			sk := el.GetSku()

			if sk.Metadata.Description.IsEmpty() {
				return
			}

			sk.Metadata.ResetTags()

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op CommitOrganizeFile) Run(
	u *env.Env,
	a, b *organize_text.Text,
	original sku.ExternalLikeSet,
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
		func(changed sku.ExternalLike) (err error) {
			// TODO switch to external
			if err = u.GetStore().CreateOrUpdateFromTransacted(
				changed.GetSku(),
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
