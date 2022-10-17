package commands

import (
	"bufio"
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/golf/alfred"
	"github.com/friedenberg/zit/src/hotel/zettel_verzeichnisse"
	"github.com/friedenberg/zit/src/mike/umwelt"
)

type CatAlfred struct {
	Type gattung.Gattung
	Command
}

func init() {
	registerCommand(
		"cat-alfred",
		func(f *flag.FlagSet) Command {
			c := &CatAlfred{
				Type: gattung.Unknown,
			}

			c.Command = c

			f.Var(&c.Type, "type", "ObjekteType")

			return c
		},
	)
}

func (c CatAlfred) Run(u *umwelt.Umwelt, args ...string) (err error) {
	//this command does its own error handling
	defer func() {
		err = nil
	}()

	wo := bufio.NewWriter(u.Out())
	defer wo.Flush()

	var aw *alfred.Writer

	if aw, err = alfred.New(wo, u.StoreObjekten()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.PanicIfError(aw.Close)

	defer func() {
		aw.WriteError(err)
	}()

	switch c.Type {
	case gattung.Etikett:
		var ea []etikett.Etikett

		if ea, err = u.StoreObjekten().Etiketten(); err != nil {
			err = errors.Wrap(err)
			return
		}

		for _, e := range ea {
			aw.WriteEtikett(e)
		}

	case gattung.Akte:
		fallthrough

	case gattung.Zettel:
		fallthrough

	case gattung.Hinweis:
		wk := zettel_verzeichnisse.MakeWriterKonfig(u.Konfig())

		if err = u.StoreObjekten().ReadAllSchwanzenVerzeichnisse(wk, aw); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = errors.Errorf("unsupported objekte type: %s", c.Type)
		return
	}

	return
}
