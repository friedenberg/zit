package command

import (
	"flag"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
)

type Request struct {
	errors.Context
	config_mutable_cli.Config
	*flag.FlagSet
	*Args
}

type consumedArg struct {
	name, value string
}

func (arg consumedArg) String() string {
	if arg.name == "" {
		return fmt.Sprintf("%q", arg.value)
	} else {
		return fmt.Sprintf("%s:%q", arg.name, arg.value)
	}
}

type Args struct {
	errors.Context
	args []string
	argi int

	consumed []consumedArg
}

func MakeRequest(
	ctx errors.Context,
	config config_mutable_cli.Config,
	flagSet *flag.FlagSet,
) Request {
	return Request{
		Context: ctx,
		Config:  config,
		FlagSet: flagSet,
		Args: &Args{
			Context: ctx,
			args:    flagSet.Args(),
		},
	}
}

func (req *Args) PeekArgs() []string {
	args := req.args[req.argi:]
	return args
}

func (req *Args) PopArgs() []string {
	args := req.PeekArgs()

	for _, arg := range args {
		req.consumed = append(req.consumed, consumedArg{value: arg})
	}

	req.argi += len(args)
	return args
}

func (req *Args) RemainingArgCount() int {
	return len(req.args[req.argi:])
}

func (req *Args) PopArg(name string) string {
	if req.RemainingArgCount() == 0 {
		req.CancelWithBadRequestf(
			"expected positional argument (%d) %s, but only received %q",
			req.argi+1,
			name,
			req.consumed,
		)
	}

	value := req.args[req.argi]
	req.consumed = append(req.consumed, consumedArg{name: name, value: value})
	req.argi++
	return value
}

func (req *Args) AssertNoMoreArgs() {
	if req.RemainingArgCount() > 0 {
		req.CancelWithBadRequestf(
			"expected no more arguments, but have %q",
			req.PopArgs(),
		)
	}
}

func (req *Args) LastArg() (arg string, ok bool) {
	if req.RemainingArgCount() > 0 {
		ok = true
		arg = req.PopArgs()[req.RemainingArgCount()-1]
	}

	return
}
