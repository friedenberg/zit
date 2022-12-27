package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/india/konfig"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/oscar/umwelt"
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
	sockPath := args[0]

	var s *remote_messages.Stage

	if s, err = remote_messages.MakeStageReceiver(u, sockPath); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, s.Close)

	var theirCliKonfig konfig.Cli

	if theirCliKonfig, err = remote_messages.PerformHi(s, u); err != nil {
		err = errors.Wrap(err)
		return
	}

	u.KonfigPtr().SetCli(theirCliKonfig)

	if err = s.Send("listening"); err != nil {
		err = errors.Wrap(err)
		return
	}

	{
		var msg string

		if err = s.Receive(&msg); err != nil {
			err = errors.Wrap(err)
			return
		}

		errors.Err().Print(msg)
	}

	{
		var msg string

		if err = s.Receive(&msg); err != nil {
			err = errors.Wrap(err)
			return
		}

		errors.Err().Print(msg)
	}

	var ids id_set.Set

	if err = s.Receive(&ids); err != nil {
		err = errors.Wrap(err)
		return
	}

	// toSend := zettel.MakeMutableSetUnique(0)

	if err = u.StoreObjekten().Zettel().ReadAllSchwanzenVerzeichnisse(
		collections.MakeChain(
			zettel.WriterIds{
				Filter: id_set.Filter{
					// AllowEmpty: true,
					Set: ids,
				},
			}.WriteZettelVerzeichnisse,
			func(z *zettel.Transacted) (err error) {
				return s.Send(z)
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// if err = s.Send(toSend.Elements()); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	return
}
