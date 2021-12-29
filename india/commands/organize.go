package commands

import (
	"flag"
	"fmt"
	"os"
)

type Organize struct {
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
	if c.rootEtiketten, err = c.getEtikettenFromArgs(args); err != nil {
		err = _Error(err)
		return
	}

	var zettels map[string]_NamedZettel

	if zettels, err = zs.Query(_NamedZettelFilterEtikettSet(c.rootEtiketten)); err != nil {
		err = _Error(err)
		return
	}

	var ot _OrganizeText

	options := _OrganizeTextOptions{
		Grouper:       c,
		Sorter:        c,
		RootEtiketten: c.rootEtiketten,
	}

	if ot, err = _OrganizeTextNew(options, zettels); err != nil {
		err = _Error(err)
		return
	}

	var p string

	err = func() (err error) {
		var f *os.File

		if f, err = _TempFileWithPattern("*.md"); err != nil {
			err = _Error(err)
			return
		}

		defer _Close(f)

		if _, err = ot.WriteTo(f); err != nil {
			err = _Error(err)
			return
		}

		p = f.Name()

		return
	}()
	//TODO remove temp

	if err != nil {
		err = _Error(err)
		return
	}

	var ot2 _OrganizeText

	for {
		if ot2, err = c.readChanges(p); err == nil {
			break
		}

		if c.handleReadChangesError(err) {
			continue
		} else {
			_Errf("aborting organize\n")
			return
		}
	}

	changes := ot.ChangesFrom(ot2)

	if err = c.commitChanges(zs, changes); err != nil {
		err = _Error(err)
		return
	}

	return
}

func (c Organize) GroupZettel(z _NamedZettel) (ess []_EtikettSet) {
	var set _EtikettSet

	if c.GroupBy.Len() > 0 {
		set = z.Zettel.Etiketten.IntersectPrefixes(c.GroupBy)
	} else {
		set = z.Zettel.Etiketten
	}

	set = set.Subtract(c.rootEtiketten)

	if c.GroupByUnique {
		ess = append(ess, set)
	} else {
		for _, e := range set {
			ns := _EtikettNewSet()
			ns.Add(e)
			ess = append(ess, ns)
		}
	}

	return ess
}

func (c Organize) SortGroups(a, b _EtikettSet) bool {
	return a.String() < b.String()
}

func (c Organize) SortZettels(a, b _NamedZettel) bool {
	return a.Hinweis.String() < b.Hinweis.String()
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

func (co Organize) readChanges(p string) (ot _OrganizeText, err error) {
	vimArgs := []string{
		"-c",
		"set ft=zit.organize",
		//TODO find a better solution for this
		"-c",
		"source ~/.vim/syntax/zit.organize.vim",
	}

	if err = _OpenVimWithArgs(vimArgs, p); err != nil {
		err = _Error(err)
		return
	}

	var f *os.File

	if f, err = _Open(p); err != nil {
		err = _Error(err)
		return
	}

	defer _Close(f)

	ot = _OrganizeTextNewEmpty()

	if _, err = ot.ReadFrom(f); err != nil {
		err = _Error(err)
		return
	}

	return
}

func (co Organize) handleReadChangesError(err error) (tryAgain bool) {
	var errorRead _OrganizeTextErrorRead

	if !_ErrorAs(err, &errorRead) {
		_Errf("unrecoverable organize read failure: %w", err)
		tryAgain = false
		return
	}

	_Errf("reading changes failed: %q\n", err)
	_Errf("would you like to edit and try again? (y/*)\n")

	var answer rune
	var n int

	if n, err = fmt.Scanf("%c", &answer); err != nil {
		tryAgain = false
		_Errf("failed to read answer: %w", err)
		return
	}

	if n != 1 {
		tryAgain = false
		_Errf("failed to read at exactly 1 answer: %w", err)
		return
	}

	if answer == 'y' || answer == 'Y' {
		tryAgain = true
	}

	return
}

func (co Organize) commitChanges(zs _Zettels, changes _OrganizeChanges) (err error) {
	if len(changes.Added) == 0 && len(changes.Removed) == 0 {
		_Out("no changes")
		return
	}

	toUpdate := make(map[string]_NamedZettel)

	addOrGetToZettelToUpdate := func(hString string) (z _NamedZettel, err error) {
		var h _Hinweis

		if h, err = _MakeBlindHinweis(hString); err != nil {
			err = _Error(err)
			return
		}

		var ok bool

		if z, ok = toUpdate[h.String()]; !ok {
			if z, err = zs.Read(h); err != nil {
				err = _Error(err)
				return
			}
		}

		return
	}

	addEtikettToZettel := func(hString string, e _Etikett) (err error) {
		var z _NamedZettel

		if z, err = addOrGetToZettelToUpdate(hString); err != nil {
			err = _Error(err)
			return
		}

		z.Zettel.Etiketten.Add(e)
		toUpdate[z.Hinweis.String()] = z

		_Outf("Added etikett '%s' to zettel '%s'\n", e, z.Hinweis)

		return
	}

	removeEtikettFromZettel := func(hString string, e _Etikett) (err error) {
		var z _NamedZettel

		if z, err = addOrGetToZettelToUpdate(hString); err != nil {
			err = _Error(err)
			return
		}

		z.Zettel.Etiketten.Remove(e)
		toUpdate[z.Hinweis.String()] = z

		_Outf("Removed etikett '%s' from zettel '%s'\n", e, z.Hinweis)

		return
	}

	for _, c := range changes.Added {
		var e _Etikett

		if err = e.Set(c.Etikett); err != nil {
			err = _Error(err)
			return
		}

		if err = addEtikettToZettel(c.Hinweis, e); err != nil {
			err = _Error(err)
			return
		}
	}

	for _, c := range changes.Removed {
		var e _Etikett

		if err = e.Set(c.Etikett); err != nil {
			err = _Error(err)
			return
		}

		if err = removeEtikettFromZettel(c.Hinweis, e); err != nil {
			err = _Error(err)
			return
		}
	}

	for _, z := range toUpdate {
		if zs.Konfig().DryRun {
			_Outf("[%s] (would update)\n", z.Hinweis)
			continue
		}

		if _, err = zs.Update(z); err != nil {
			_Errf("failed to update zettel: %w", err)
		}
	}

	return
}
