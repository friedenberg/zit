package user_ops

import (
	"os"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/golf/organize_text"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type CreateOrganizeFile struct {
	Umwelt        *umwelt.Umwelt
	RootEtiketten etikett.Set
	GroupBy       etikett.Set
	GroupByUnique bool
}

type CreateOrgaanizeFileResults struct {
	Path string
	Text organize_text.Text
}

func (c CreateOrganizeFile) Run(zettels map[string]stored_zettel.Named) (results CreateOrgaanizeFileResults, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	options := organize_text.Options{
		Grouper:       c,
		Sorter:        c,
		RootEtiketten: c.RootEtiketten,
	}

	if results.Text, err = organize_text.New(options, zettels); err != nil {
		err = _Error(err)
		return
	}

	err = func() (err error) {
		var f *os.File

		if f, err = open_file_guard.TempFileWithPattern("*.md"); err != nil {
			err = _Error(err)
			return
		}

		defer open_file_guard.Close(f)

		if _, err = results.Text.WriteTo(f); err != nil {
			err = _Error(err)
			return
		}

		results.Path = f.Name()

		return
	}()
	//TODO remove temp

	if err != nil {
		err = _Error(err)
		return
	}

	return
}

func (c CreateOrganizeFile) GroupZettel(z _NamedZettel) (ess []etikett.Set) {
	var set etikett.Set

	if c.GroupBy.Len() > 0 {
		set = z.Zettel.Etiketten.IntersectPrefixes(c.GroupBy)
	} else {
		set = z.Zettel.Etiketten
	}

	set = set.Subtract(c.RootEtiketten)

	if c.GroupByUnique {
		ess = append(ess, set)
	} else {
		for _, e := range set {
			ns := etikett.NewSet()
			ns.Add(e)
			ess = append(ess, ns)
		}
	}

	return ess
}

func (c CreateOrganizeFile) SortGroups(a, b etikett.Set) bool {
	return a.String() < b.String()
}

func (c CreateOrganizeFile) SortZettels(a, b _NamedZettel) bool {
	return a.Hinweis.String() < b.Hinweis.String()
}

func (c CreateOrganizeFile) getEtikettenFromArgs(args []string) (es etikett.Set, err error) {
	es = etikett.NewSet()

	for _, s := range args {
		if err = es.AddString(s); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}
