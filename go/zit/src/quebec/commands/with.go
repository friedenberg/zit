package commands

import (
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

type WithEnv interface {
	RunWithEnv(*env.Env, ...string)
}

type WithArchive interface {
	Run(repo.Archive, ...string)
}

type WithWorkingCopy interface {
	RunWithWorkingCopy(repo.WorkingCopy, ...string)
}

type WithLocalWorkingCopy interface {
	RunWithLocalWorkingCopy(*local_working_copy.Repo, ...string)
}

type WithBlobStore interface {
	RunWithBlobStore(command_components.BlobStoreWithEnv, ...string)
}

type WithQuery interface {
	RunWithQuery(store *local_working_copy.Repo, ids *query.Group)
}

type WithQueryAndBuilderOptions interface {
	query.BuilderOptionGetter
	WithQuery
}
