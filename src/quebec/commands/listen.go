package commands

import (
	"flag"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/remote_conn"
	"github.com/friedenberg/zit/src/quebec/remote_pull"
)

type Listen struct {
}

func init() {
	registerCommand(
		"listen",
		func(f *flag.FlagSet) Command {
			c := &Listen{}

			return c
		},
	)
}

func (c Listen) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) == 0 {
		err = errors.Normalf("must specify command to listen for")
		return
	}

	command := args[0]
	var l remote_conn.Listener

	switch strings.ToLower(strings.TrimSpace(command)) {
	case "pull":
		if l, err = remote_pull.MakeServer(u); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = errors.Normalf("unsupported command: %q", command)
		return
	}

	if err = l.Listen(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
