package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/india/store_with_lock"
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

			return c
		},
	)
}

func (c *Organize) Run(u _Umwelt, args ...string) (err error) {
	createOrganizeFileOp := user_ops.CreateOrganizeFile{
		Umwelt:        u,
		GroupBy:       c.GroupBy,
		GroupByUnique: c.GroupByUnique,
	}

	if createOrganizeFileOp.RootEtiketten, err = c.getEtikettenFromArgs(args); err != nil {
		err = _Error(err)
		return
	}

	var zettels map[string]_NamedZettel

	if zettels, err = c.getZettels(u, createOrganizeFileOp.RootEtiketten); err != nil {
		err = _Error(err)
		return
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

func (c Organize) getZettels(u _Umwelt, rootEtiketten _EtikettSet) (zettels map[string]_NamedZettel, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(u); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	if zettels, err = store.Zettels().Query(_NamedZettelFilterEtikettSet(rootEtiketten)); err != nil {
		err = _Error(err)
		return
	}

	return
}
