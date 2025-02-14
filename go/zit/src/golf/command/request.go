package command

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
)

type Request struct {
	errors.Context
	config_mutable_cli.Config
	*flag.FlagSet

	args []string
	argi *int
}

func MakeRequest(
	ctx errors.Context,
	config config_mutable_cli.Config,
	flagSet *flag.FlagSet,
) Request {
	argi := 0

	return Request{
		Context: ctx,
		Config:  config,
		FlagSet: flagSet,
		args:    flagSet.Args(),
		argi:    &argi,
	}
}

func (req Request) Args() []string {
	return req.args[*req.argi:]
}

func (req Request) Argc() int {
	return len(req.Args())
}

func (req Request) Argv(argName string) string {
	if req.Argc() == 0 {
		req.CancelWithBadRequestf(
			"expected positional argument (%d) %s, but only received %q",
			argName,
			req.args,
		)
	}

	arg := req.args[*req.argi]
	*req.argi++
	return arg
}

func (req Request) AssertNoMoreArgs() {
	if req.Argc() > 0 {
		req.CancelWithBadRequestf(
			"expected no more arguments, but have %q",
			req.Args(),
		)
	}
}

func (req Request) LastArg() (arg string, ok bool) {
	if req.Argc() > 0 {
		ok = true
		arg = req.Args()[req.Argc()-1]
	}

	return
}
