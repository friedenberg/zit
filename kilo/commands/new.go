package commands

import (
	"flag"

	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
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

			return commandWithZettels{c}
		},
	)
}

func (c New) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	f := _ZettelFormatsText{}

	if u.Konfig.DryRun && len(args) == 0 {
		_Errf("when -dry-run is set, paths to existing zettels must be provided")
		return
	}

	if len(args) == 0 {
		if err = c.writeEmptyAndOpen(u, zs, f); err != nil {
			err = _Error(err)
			return
		}
	} else {
		newOp := user_ops.CreateFromPaths{
			Umwelt: u,
			Store:  zs,
			Format: f,
			Filter: c.Filter,
		}

		if _, err = newOp.Run(args...); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}

func (c New) writeEmptyAndOpen(u _Umwelt, zs _Zettels, format zettel.Format) (err error) {
	newOp := user_ops.WriteEmptyZettel{
		Umwelt: u,
		Store:  zs,
		Format: format,
	}

	var results user_ops.WriteEmptyZettelResults

	if results, err = newOp.Run(); err != nil {
		err = _Error(err)
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: []string{
			//TODO move to builder
			`call cursor(2, 3)`,
			`startinsert!`,
			"set ft=zit.zettel",
			"source ~/.vim/syntax/zit.zettel.vim",
		},
	}

	// var openVimResults user_ops.OpenVimResults

	if _, err = openVimOp.Run(results.Zettel.Path); err != nil {
		err = _Error(err)
		return
	}

	options := _ZettelsCheckinOptions{
		IncludeAkte: true,
		Format:      _ZettelFormatsText{},
	}

	var external map[hinweis.Hinweis]stored_zettel.External

	if external, err = zs.ReadExternal(options, results.Zettel.Path); err != nil {
		err = _Error(err)
		return
	}

	for _, z := range external {
		// var named _NamedZettel

		if _, err = zs.CreateWithHinweis(z.Zettel, z.Hinweis); err != nil {
			err = _Error(err)
			return
		}
	}

	// checkinOp := user_ops.Checkin{
	// 	Umwelt: u,
	// 	Store:  zs,
	// }

	// if _, err = checkinOp.Run(results.Zettel.Path); err != nil {
	// 	err = _Error(err)
	// 	return
	// }

	return
}
