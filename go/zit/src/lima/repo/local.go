package repo

import "code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"

type LocalRepo interface {
	Repo
	GetRepoLayout() env_repo.Env
}

type LocalWorkingCopy interface {
	WorkingCopy
	LocalRepo
}
