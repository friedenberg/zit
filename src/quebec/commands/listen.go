package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
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
	var s *remote_messages.StageSoldier
	var theirHiMessage remote_messages.MessageHiCommander

	if s, theirHiMessage, err = remote_messages.MakeStageSoldier(u); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO-P1 fix verbose logging on remotes
	u.KonfigPtr().SetCliFromCommander(theirHiMessage.CliKonfig)

	if err = u.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Print(u.Standort().Cwd())

	errors.Log().Printf("setting konfig")

	s.RegisterHandlerWithUmwelt(
		remote_messages.DialogueTypePull,
		u,
		c.handleDialoguePull,
	)

	s.RegisterHandlerWithUmwelt(
		remote_messages.DialogueTypePullObjekten,
		u,
		c.handleDialoguePullObjekten,
	)

	s.RegisterHandlerWithUmwelt(
		remote_messages.DialogueTypePullAkte,
		u,
		c.handleDialoguePullAkte,
	)

	if err = s.Listen(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Listen) handleDialoguePullAkte(
	u *umwelt.Umwelt,
	d remote_messages.Dialogue,
) (err error) {
	var sh sha.Sha

	errors.Log().Print("waiting to receive sha")

	if err = d.Receive(&sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, d.Close)

	errors.Log().Printf("did receive sha: %s", sh)

	var ar sha.ReadCloser

	if ar, err = u.StoreObjekten().AkteReader(sh); err != nil {
		errors.Log().Printf("got error on akte reader: %s", err)
		if errors.IsNotExist(err) {
			err = nil
			return
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	var n int64

	if n, err = io.Copy(d, ar); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("copied %d bytes", n)

	return
}

func (c Listen) handleDialoguePull(
	u *umwelt.Umwelt,
	d remote_messages.Dialogue,
) (err error) {
	var filter id_set.Filter

	if err = d.Receive(&filter); err != nil {
		err = errors.Wrap(err)
		return
	}

	t := transaktion.MakeTransaktion(ts.Now())

	if err = u.StoreObjekten().Zettel().ReadAllSchwanzenVerzeichnisse(
		collections.MakeChain(
			zettel.WriterIds{Filter: filter}.WriteZettelVerzeichnisse,
			func(z *zettel.Transacted) (err error) {
				t.Skus.Add2(&z.Sku)
				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = d.Send(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Listen) handleDialoguePullObjekten(
	u *umwelt.Umwelt,
	d remote_messages.Dialogue,
) (err error) {
	skus := sku.MakeMutableSet()

	errors.Log().Print("waiting to receive skus")

	if err = d.Receive(&skus); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("did receive skus: %d", skus.Len())

	if err = skus.Each(
		func(s *sku.Sku) (err error) {
			//TODO-P1 support any transacted objekte
			if s.Gattung != gattung.Zettel {
				return
			}

			h := s.Id.(hinweis.Hinweis)

			var zt *zettel.Transacted

			if zt, err = u.StoreObjekten().Zettel().ReadHinweisSchwanzen(h); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = d.Send(zt); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = d.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
