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

func (req Request) LastArg() (arg string, ok bool) {
	argc := len(req.Args())

	if argc > 0 {
		ok = true
		arg = req.Args()[argc-1]
	}

	return
}
