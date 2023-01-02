package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/user_ops"
	"github.com/friedenberg/zit/src/remote_messages"
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
	var s *remote_messages.StageSoldier

	if s, err = remote_messages.MakeStageSoldier(u); err != nil {
		err = errors.Wrap(err)
		return
	}

	pullOrPushOp := user_ops.MakeRemoteMessagesPullOrPush(u)
	pullOrPushOp.AddToSoldierStage(s)

	if err = s.Listen(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
