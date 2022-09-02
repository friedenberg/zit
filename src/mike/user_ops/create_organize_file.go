package user_ops

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/delta/umwelt"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/organize_text"
	"github.com/friedenberg/zit/src/lima/store_with_lock"
)

type CreateOrganizeFile struct {
	Umwelt *umwelt.Umwelt
	organize_text.Options
}

type CreateOrganizeFileResults struct {
	Text organize_text.Text
}

func (c CreateOrganizeFile) RunAndWrite(zettels zettel_transacted.Set, w io.WriteCloser) (results CreateOrganizeFileResults, err error) {
	results, err = c.Run(zettels)

	defer errors.PanicIfError(w.Close)

	if _, err = results.Text.WriteTo(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c CreateOrganizeFile) Run(zettels zettel_transacted.Set) (results CreateOrganizeFileResults, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	if results.Text, err = organize_text.New(c.Options); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// func (c CreateOrganizeFile) GroupZettel(z stored_zettel.Named) (ess []etikett.Set) {
// 	var set etikett.Set

// 	if c.GroupBy.Len() > 0 {
// 		set = z.Zettel.Etiketten.IntersectPrefixes(c.GroupBy.ToSet())
// 	} else {
// 		set = z.Zettel.Etiketten
// 	}

// 	set = set.Subtract(c.RootEtiketten)

// 	if false /*c.GroupByUnique*/ {
// 		ess = append(ess, set)
// 	} else if set.Len() > 0 {
// 		for _, e := range set {
// 			ns := etikett.MakeSet()
// 			ns.Add(e)
// 			ess = append(ess, ns)
// 		}
// 	} else {
// 		// if the zettel has no etiketten, add an empty set
// 		ess = append(ess, etikett.MakeSet())
// 	}

// 	return ess
// }

// func (c CreateOrganizeFile) SortGroups(a, b etikett.Set) bool {
// 	return a.String() < b.String()
// }

// func (c CreateOrganizeFile) SortZettels(a, b stored_zettel.Named) bool {
// 	return a.Hinweis.String() < b.Hinweis.String()
// }

func (c CreateOrganizeFile) getEtikettenFromArgs(args []string) (es etikett.Set, err error) {
	es = etikett.MakeSet()

	for _, s := range args {
		if err = es.AddString(s); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
