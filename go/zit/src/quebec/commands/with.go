package commands

import (
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type WithWorkingCopy interface {
	Run(repo.WorkingCopy, ...string)
}

type WithLocalWorkingCopy interface {
	Run(*local_working_copy.Repo, ...string)
}

type WithQuery interface {
	Run(store *local_working_copy.Repo, ids *query.Group)
}

type WithQueryAndBuilderOptions interface {
	query.BuilderOptionGetter
	WithQuery
}
