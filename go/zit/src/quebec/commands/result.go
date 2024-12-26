package commands

import "code.linenisgreat.com/zit/go/zit/src/november/repo_local"

type Result struct {
	Success  bool
	ExitCode int
	Error    error
}

type commandWithResult struct {
	CommandWithRepo
}

func (cwr commandWithResult) RunWithRepo(u *repo_local.Repo, args ...string) {
	cwr.CommandWithRepo.RunWithRepo(u, args...)
}
