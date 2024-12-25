package commands

import "code.linenisgreat.com/zit/go/zit/src/november/repo_local"

type Result struct {
	Success  bool
	ExitCode int
	Error    error
}

type commandWithResult struct {
	Command
}

func (cwr commandWithResult) Run(u *repo_local.Repo, args ...string) {
	if err := cwr.Command.Run(u, args...); err != nil {
		u.Context.Cancel(err)
		return
	}
}
