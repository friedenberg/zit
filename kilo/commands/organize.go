package commands

import (
	"flag"

	"github.com/friedenberg/zit/juliett/user_ops"
)

type Organize struct {
	Hinweisen     bool //TODO support
	rootEtiketten _EtikettSet
	GroupBy       _EtikettSet
	GroupByUnique bool
}

func init() {
	registerCommand(
		"organize",
		func(f *flag.FlagSet) Command {
			c := &Organize{
				GroupBy: _EtikettNewSet(),
			}

			f.BoolVar(&c.GroupByUnique, "group-by-unique", false, "group by all unique combinations of etiketten")
			f.Var(&c.GroupBy, "group-by", "etikett prefixes to group zettels")

			return commandWithZettels{c}
		},
	)
}

func (c *Organize) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	var zettels map[string]_NamedZettel

	createOrganizeFileOp := user_ops.CreateOrganizeFile{
		Umwelt:        u,
		GroupBy:       c.GroupBy,
		GroupByUnique: c.GroupByUnique,
	}

	if c.Hinweisen {
		//TODO add RootEtiketten
		zettels = make(map[string]_NamedZettel)

		for _, arg := range args {
			var h _Hinweis
			if h, err = _MakeBlindHinweis(arg); err != nil {
				err = _Error(err)
				return
			}

			var named _NamedZettel

			if named, err = zs.Read(h); err != nil {
				err = _Error(err)
				return
			}

			zettels[h.String()] = named
		}
	} else {
		if createOrganizeFileOp.RootEtiketten, err = c.getEtikettenFromArgs(args); err != nil {
			err = _Error(err)
			return
		}

		if zettels, err = zs.Query(_NamedZettelFilterEtikettSet(c.rootEtiketten)); err != nil {
			err = _Error(err)
			return
		}
	}

	var createOrganizeFileResults user_ops.CreateOrgaanizeFileResults

	if createOrganizeFileResults, err = createOrganizeFileOp.Run(zettels); err != nil {
		err = _Error(err)
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: []string{
			"set ft=zit.organize",
			//TODO find a better solution for this
			"source ~/.vim/syntax/zit.organize.vim",
		},
	}

	if _, err = openVimOp.Run(createOrganizeFileResults.Path); err != nil {
		err = _Error(err)
		return
	}

	var ot2 _OrganizeText

	readOrganizeTextOp := user_ops.ReadOrganizeFile{}

	if ot2, err = readOrganizeTextOp.Run(createOrganizeFileResults.Path); err != nil {
		err = _Error(err)
		return
	}

	commitOrganizeTextOp := user_ops.CommitOrganizeFile{
		Umwelt: u,
	}

	if _, err = commitOrganizeTextOp.Run(createOrganizeFileResults.Text, ot2); err != nil {
		err = _Error(err)
		return
	}

	return
}

func (c Organize) getEtikettenFromArgs(args []string) (es _EtikettSet, err error) {
	es = _EtikettNewSet()

	for _, s := range args {
		if err = es.AddString(s); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}
