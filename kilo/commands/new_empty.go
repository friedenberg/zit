package commands

import (
	"flag"
	"log"

	"github.com/friedenberg/zit/juliett/user_ops"
)

type NewEmpty struct {
}

func init() {
	registerCommand(
		"new-empty",
		func(f *flag.FlagSet) Command {
			c := &NewEmpty{}

			return c
		},
	)
}

func (c NewEmpty) Run(u _Umwelt, _ ...string) (err error) {
	f := _ZettelFormatsText{}

	emptyOp := user_ops.WriteEmptyZettel{
		Umwelt: u,
		Format: f,
	}

	var results user_ops.WriteEmptyZettelResults

	if results, err = emptyOp.Run(); err != nil {
		err = _Error(err)
		return
	}

	log.Print(results)
	_Outf("%s\n", results.Zettel.Hinweis)

	return
}
