package user_ops

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/organize_text"
	"github.com/friedenberg/zit/src/kilo/umwelt"
)

type CreateOrganizeFile struct {
	*umwelt.Umwelt
	organize_text.Options
}

func (c CreateOrganizeFile) RunAndWrite(
	zettels zettel_transacted.Set,
	w io.WriteCloser,
) (results *organize_text.Text, err error) {
	defer errors.PanicIfError(w.Close)

	if results, err = c.Run(zettels); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = results.WriteTo(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c CreateOrganizeFile) Run(zettels zettel_transacted.Set) (results *organize_text.Text, err error) {
	if results, err = organize_text.New(c.Options); err != nil {
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
