package commands

import "code.linenisgreat.com/zit/go/zit/src/november/env"

type Result struct {
	Success  bool
	ExitCode int
	Error    error
}

type commandWithResult struct {
	Command
}

func (cwr commandWithResult) Run(u *env.Env, args ...string) (result Result) {
	result.Error = cwr.Command.Run(u, args...)
	return
}
