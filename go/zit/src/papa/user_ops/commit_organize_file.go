package user_ops

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/kilo/organize_text"
	"code.linenisgreat.com/zit/src/lima/changes"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type CommitOrganizeFile struct {
	*umwelt.Umwelt
	OutputJSON bool
}

type CommitOrganizeFileResults = changes.Changes2

func (c CommitOrganizeFile) ApplyToText(
	u *umwelt.Umwelt,
	t *organize_text.Text,
) (err error) {
	if u.Konfig().PrintOptions.PrintEtikettenAlways {
		return
	}

	if err = t.Transacted.Each(
		func(sk *sku.Transacted) (err error) {
			if sk.Metadatei.Bezeichnung.IsEmpty() {
				return
			}

			sk.Metadatei.GetEtikettenMutable().Reset()

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

	if cs, err = changes.ChangesFrom2(
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

	if err = u.GetStore().UpdateManyMetadatei(cs.B); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
