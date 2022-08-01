package commands

import (
	"bufio"
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type CatAlfred struct {
	Type _Type
	Command
}

func init() {
	registerCommand(
		"cat-alfred",
		func(f *flag.FlagSet) Command {
			c := &CatAlfred{
				Type: _TypeUnknown,
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

	var aw _AlfredWriter

	if aw, err = _AlfredNewWriter(store.Out); err != nil {
		err = errors.Error(err)
		return
	}

	defer stdprinter.PanicIfError(aw.Close)

	defer func() {
		aw.WriteError(err)
	}()

	switch c.Type {
	case _TypeEtikett:
		var ea []_Etikett

		if ea, err = store.Etiketten().All(); err != nil {
			err = errors.Error(err)
			return
		}

		for _, e := range ea {
			aw.WriteEtikett(e)
		}

	case _TypeZettel:

		var all map[string]_NamedZettel

		if all, err = store.Zettels().All(); err != nil {
			err = errors.Error(err)
			return
		}

		for _, z := range all {
			aw.WriteZettel(z)
		}

	case _TypeAkte:

	case _TypeHinweis:

		var all map[string]_NamedZettel

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
