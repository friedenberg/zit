package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
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

func (op CommitOrganizeFile) RunTraditionalCommit(
	u *env.Env,
	a, b *organize_text.Text,
	original sku.ExternalLikeSet,
	qg *query.Group,
) (cs CommitOrganizeFileResults, err error) {
	if cs, err = u.CommitOrganizeResults(
		organize_text.OrganizeResults{
			Before:     a,
			After:      b,
			Original:   original,
			QueryGroup: qg,
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}
	if cs, err = organize_text.ChangesFrom(
		u.GetConfig().PrintOptions,
		a, b, original,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
