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
	Args
}

type Args struct {
	errors.Context
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
		Args: Args{
			Context: ctx,
			args:    flagSet.Args(),
			argi:    &argi,
		},
	}
}

func (req Args) PeekArgs() []string {
	args := req.args[*req.argi:]
	return args
}

func (req Args) PopArgs() []string {
	args := req.PeekArgs()
	*req.argi += len(args)
	return args
}

func (req Args) RemainingArgCount() int {
	return len(req.args[*req.argi:])
}

func (req Args) PopArg(argName string) string {
	if req.RemainingArgCount() == 0 {
		req.CancelWithBadRequestf(
			"expected positional argument (%d) %s, but only received %q",
			*req.argi,
			argName,
			req.args,
		)
	}

	arg := req.args[*req.argi]
	*req.argi++
	return arg
}

func (req Args) AssertNoMoreArgs() {
	if req.RemainingArgCount() > 0 {
		req.CancelWithBadRequestf(
			"expected no more arguments, but have %q",
			req.PopArgs(),
		)
	}
}

func (req Args) LastArg() (arg string, ok bool) {
	if req.RemainingArgCount() > 0 {
		ok = true
		arg = req.PopArgs()[req.RemainingArgCount()-1]
	}

	return
}
