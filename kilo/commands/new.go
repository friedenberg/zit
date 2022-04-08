package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/juliett/user_ops"
)

type New struct {
	Filter _ScriptValue
}

func init() {
	registerCommand(
		"new",
		func(f *flag.FlagSet) Command {
			c := &New{}

			f.Var(&c.Filter, "filter", "a script to run for each file to transform it the standard zettel format")

			return c
		},
	)
}

func (c New) Run(u _Umwelt, args ...string) (err error) {
	f := _ZettelFormatsText{}

	if u.Konfig.DryRun && len(args) == 0 {
		_Errf("when -dry-run is set, paths to existing zettels must be provided")
		return
	}

	newOp := user_ops.CreateFromPaths{
		Umwelt: u,
		Format: f,
		Filter: c.Filter,
	}

	if len(args) == 0 {
		emptyOp := user_ops.WriteEmptyZettel{
			Umwelt: u,
			Format: f,
		}

		var results user_ops.WriteEmptyZettelResults

		if results, err = emptyOp.Run(); err != nil {
			err = _Error(err)
			return
		}

		openVimOp := user_ops.OpenVim{
			Options: vim_cli_options_builder.New().
				WithCursorLocation(2, 3).
				WithInsertMode().
				WithFileType("zit.zettel").
				WithSourcedFile("~/.vim/syntax/zit.zettel.vim").
				Build(),
		}

		if _, err = openVimOp.Run(results.Zettel.Path); err != nil {
			err = _Error(err)
			return
		}

		newOp.ReadHinweisFromPath = true
		args = []string{results.Zettel.Path}
	}

	if _, err = newOp.Run(args...); err != nil {
		err = _Error(err)
		return
	}

	return
}
