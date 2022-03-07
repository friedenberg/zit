package commands

import (
	"flag"
	"io"
	"os"
	"path"
)

type Add struct {
	Etiketten _EtikettSet
	Delete    bool
	Organize  bool
	//TODO
	// Edit      bool
}

func init() {
	registerCommand(
		"add",
		func(f *flag.FlagSet) Command {
			c := &Add{
				Etiketten: _EtikettNewSet(),
			}

			f.Var(&c.Etiketten, "etiketten", "to add to the created zettels")
			f.BoolVar(&c.Delete, "delete", false, "delete the zettel and akte after successful checkin")
			f.BoolVar(&c.Organize, "organize", false, "")
			// f.BoolVar(&c.Edit, "edit", false, "")

			return commandWithZettels{c}
		},
	)
}

func (c Add) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	added := make([]_NamedZettel, len(args))
	hinweisen := make([]string, len(args))

	for i, arg := range args {
		var z _Zettel

		if z, err = c.zettelForAkte(u, zs, arg); err != nil {
			err = _Error(err)
			return
		}

		var named _NamedZettel

		if named, err = zs.Create(z); err != nil {
			err = _Error(err)
			return
		}

		added[i] = named
		hinweisen[i] = named.Hinweis.String()

		if c.Delete {
			if err = os.Remove(arg); err != nil {
				err = _Error(err)
				return
			}

			_Errf("[%s] (deleted)\n", arg)
		}

		_Outf("[%s %s] (created)\n", named.Hinweis, named.Sha)
	}

	//TODO move to user ops
	if c.Organize {
		c1 := &Organize{
			Hinweisen:     true,
			GroupBy:       _EtikettNewSet(),
			GroupByUnique: true,
		}

		if err = c1.RunWithZettels(u, zs, hinweisen...); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}

func (c Add) zettelForAkte(u _Umwelt, zs _Zettels, aktePath string) (z _Zettel, err error) {
	z.Etiketten = c.Etiketten

	var akteWriter _ObjekteWriter

	if akteWriter, err = zs.AkteWriter(); err != nil {
		err = _Error(err)
		return
	}

	var f *os.File

	if f, err = _Open(aktePath); err != nil {
		err = _Error(err)
		return
	}

	defer _Close(f)

	if _, err = io.Copy(akteWriter, f); err != nil {
		err = _Error(err)
		return
	}

	if err = akteWriter.Close(); err != nil {
		err = _Error(err)
		return
	}

	if err = z.Bezeichnung.Set(path.Base(aktePath)); err != nil {
		err = _Error(err)
		return
	}

	z.Akte = akteWriter.Sha()

	if err = z.AkteExt.Set(path.Ext(aktePath)); err != nil {
		err = _Error(err)
		return
	}

	return
}
