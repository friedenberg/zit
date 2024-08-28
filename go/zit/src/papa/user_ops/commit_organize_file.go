package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
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

	if err = t.Options.Transacted.Each(
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

func (op CommitOrganizeFile) RunCommit(
	u *env.Env,
	a, b *organize_text.Text,
	original sku.ExternalLikeSet,
	qg *query.Group,
	onChanged interfaces.FuncIter[sku.ExternalLike],
) (cs CommitOrganizeFileResults, err error) {
	if err = op.ApplyToText(u, a); err != nil {
		err = errors.Wrap(err)
		return
	}

	if cs, err = organize_text.ChangesFrom(a, b, original); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO move to bestandsaufnahme flush
	// if cs.GetAddedUnnamed().Len() == 0 && cs.GetAddedNamed().Len() == 0 {
	// 	errors.Err().Print("no changes")
	// 	return
	// }

	if onChanged == nil {
		onChanged = func(changed sku.ExternalLike) (err error) {
			if err = u.GetStore().CreateOrUpdate(
				changed,
				objekte_mode.Make(
					objekte_mode.ModeMergeCheckedOut,
				),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	if err = cs.Changed.Each(onChanged); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
