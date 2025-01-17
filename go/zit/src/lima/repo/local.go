package repo

import "code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"

type LocalRepo interface {
	Repo
	GetRepoLayout() repo_layout.Layout
}

type LocalWorkingCopy interface {
	WorkingCopy
	LocalRepo
}
