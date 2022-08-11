package commands

import (
	"bufio"
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/node_type"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/golf/alfred"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type CatAlfred struct {
	Type node_type.Type
	Command
}

func init() {
	registerCommand(
		"cat-alfred",
		func(f *flag.FlagSet) Command {
			c := &CatAlfred{
				Type: node_type.TypeUnknown,
			}

			c.Command = commandWithLockedStore{c}

			f.Var(&c.Type, "type", "ObjekteType")

			return c
		},
	)
}

func (c CatAlfred) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	//this command does its own error handling
	defer func() {
		err = nil
	}()

	wo := bufio.NewWriter(store.Out)
	defer wo.Flush()

	var aw alfred.Writer

	if aw, err = alfred.NewWriter(store.Out); err != nil {
		err = errors.Error(err)
		return
	}

	defer stdprinter.PanicIfError(aw.Close)

	defer func() {
		aw.WriteError(err)
	}()

	switch c.Type {
	case node_type.TypeEtikett:
		var ea []etikett.Etikett

		if ea, err = store.Zettels().Etiketten(); err != nil {
			err = errors.Error(err)
			return
		}

		for _, e := range ea {
			aw.WriteEtikett(e)
		}

	case node_type.TypeZettel:

		var all map[hinweis.Hinweis]stored_zettel.Transacted

		if all, err = store.Zettels().ZettelTails(); err != nil {
			err = errors.Error(err)
			return
		}

		for _, z := range all {
			aw.WriteZettel(z.Named)
		}

	case node_type.TypeAkte:

	case node_type.TypeHinweis:

		var all map[string]stored_zettel.Named

		if all, err = store.Zettels().All(); err != nil {
			err = errors.Error(err)
			return
		}

		for _, z := range all {
			aw.WriteZettel(z)
		}

	default:
		err = errors.Errorf("unsupported objekte type: %s", c.Type)
		return
	}

	return
}
