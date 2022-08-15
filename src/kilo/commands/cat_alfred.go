package commands

import (
	"bufio"
	"flag"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/charlie/zk_types"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	"github.com/friedenberg/zit/src/hotel/alfred"
	"github.com/friedenberg/zit/src/india/store_with_lock"
)

type CatAlfred struct {
	Type zk_types.Type
	Command
}

func init() {
	registerCommand(
		"cat-alfred",
		func(f *flag.FlagSet) Command {
			c := &CatAlfred{
				Type: zk_types.TypeUnknown,
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
	case zk_types.TypeEtikett:
		var ea []etikett.Etikett

		if ea, err = store.Zettels().Etiketten(); err != nil {
			err = errors.Error(err)
			return
		}

		for _, e := range ea {
			aw.WriteEtikett(e)
		}

	case zk_types.TypeZettel:

		var all map[hinweis.Hinweis]stored_zettel.Transacted

		if all, err = store.Zettels().ZettelTails(); err != nil {
			err = errors.Error(err)
			return
		}

		for _, z := range all {
			aw.WriteZettel(z.Named)
		}

	case zk_types.TypeAkte:

	case zk_types.TypeHinweis:

		var all map[hinweis.Hinweis]stored_zettel.Transacted

		if all, err = store.Zettels().ZettelTails(); err != nil {
			err = errors.Error(err)
			return
		}

		for _, z := range all {
			aw.WriteZettel(z.Named)
		}

	default:
		err = errors.Errorf("unsupported objekte type: %s", c.Type)
		return
	}

	return
}
