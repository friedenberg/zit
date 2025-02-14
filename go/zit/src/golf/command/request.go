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
}

func (req Request) Argc() int {
	return len(req.Args())
}

func (req Request) Argv(idx int, argName string) string {
	if req.Argc()-1 < idx {
		req.CancelWithBadRequestf(
			"expected %s at position %d, but only received %q",
			argName,
			idx,
			req.Args(),
		)
	}

	return req.Args()[idx]
}

func (req Request) LastArg() (arg string, ok bool) {
	if req.Argc() > 0 {
		ok = true
		arg = req.Args()[req.Argc()-1]
	}

	return
}
