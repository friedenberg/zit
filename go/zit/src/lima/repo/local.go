package repo

import (
	"code.linenisgreat.com/zit/go/zit/src/golf/config_immutable_io"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/kilo/env_workspace"
)

type LocalRepo interface {
	Repo
	GetEnvRepo() env_repo.Env // TODO rename to GetEnvRepo
	GetImmutableConfigPrivate() config_immutable_io.ConfigLoadedPrivate
	Lock() error
	Unlock() error
}

type LocalWorkingCopy interface {
	WorkingCopy
	LocalRepo
	GetEnvWorkspace() env_workspace.Env
}
