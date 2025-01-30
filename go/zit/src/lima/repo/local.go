package repo

import "code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"

type LocalRepo interface {
	Repo
	GetEnvRepo() env_repo.Env // TODO rename to GetEnvRepo
}

type LocalWorkingCopy interface {
	WorkingCopy
	LocalRepo
}
