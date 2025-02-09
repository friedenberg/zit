package command

import "code.linenisgreat.com/zit/go/zit/src/hotel/env_local"

type SupportsCompletion interface {
	SupportsCompletion()
}

type CommandLine struct {
	Args       []string
	InProgress string
}

func (commandLine CommandLine) LastArg() (arg string, ok bool) {
	argc := len(commandLine.Args)

	if argc > 0 {
		ok = true
		arg = commandLine.Args[argc-1]
	}

	return
}

func (commandLine CommandLine) LastCompleteArg() (arg string, ok bool) {
	argc := len(commandLine.Args)

	if commandLine.InProgress != "" {
		argc -= 1
	}

	if argc > 0 {
		ok = true
		arg = commandLine.Args[argc-1]
	}

	return
}

type Completion struct {
	Value, Description string
}

type Completer interface {
	Complete(Request, env_local.Env, CommandLine)
}
