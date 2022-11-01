package store_objekten

import (
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/india/zettel_verzeichnisse"
	"github.com/friedenberg/zit/src/juliett/store_verzeichnisse"
)

type verzeichnisseAll struct {
	*store_verzeichnisse.Zettelen
	ioFactory
}

func makeVerzeichnisseAll(
	k konfig.Konfig,
	st standort.Standort,
	iof ioFactory,
	p zettel_verzeichnisse.Pool,
) (s *verzeichnisseAll, err error) {
	s = &verzeichnisseAll{
		ioFactory: iof,
	}

	s.Zettelen, err = store_verzeichnisse.MakeZettelen(
		k,
		st.DirVerzeichnisseZettelenNeue(),
		s,
		p,
	)

	return
}

// func (s *verzeichnisseAll) add(tz *zettel_transacted.Zettel) {
// 	if err = s.Zettelen.Add(tz, tz.Named.Hinweis.String()); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	{
// 		var set zettel_transacted.Set
// 		var ok bool

// 		if set, ok = i.hinweisen[tz.Named.Hinweis]; !ok {
// 			set = zettel_transacted.MakeSetUnique(1)
// 		}

// 		set.Add(tz)
// 		i.hinweisen[tz.Named.Hinweis] = set
// 	}

// 	akteSha := tz.Named.Stored.Zettel.Akte

// 	if !akteSha.IsNull() {
// 		var set zettel_transacted.Set
// 		var ok bool

// 		if set, ok = i.akten[tz.Named.Stored.Zettel.Akte]; !ok {
// 			set = zettel_transacted.MakeSetUnique(1)
// 		}

// 		set.Add(tz)
// 		i.akten[tz.Named.Stored.Zettel.Akte] = set
// 	}

// 	bezKey := strings.ToLower(tz.Named.Stored.Zettel.Bezeichnung.String())
// 	if bezKey != "" {

// 		var set zettel_transacted.Set
// 		var ok bool

// 		if set, ok = i.bezeichnungen[bezKey]; !ok {
// 			set = zettel_transacted.MakeSetUnique(1)
// 		}

// 		set.Add(tz)
// 		i.bezeichnungen[bezKey] = set
// 	}

// 	{
// 		var set zettel_transacted.Set
// 		var ok bool

// 		if set, ok = i.typen[tz.Named.Stored.Zettel.Typ]; !ok {
// 			set = zettel_transacted.MakeSetUnique(1)
// 		}

// 		set.Add(tz)
// 		i.typen[tz.Named.Stored.Zettel.Typ] = set
// 	}

// 	return
// }

// func (i *indexZettelen) Add(tz zettel_transacted.Zettel) (err error) {
// 	if err = i.readIfNecessary(); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	i.hasChanges = true

// 	i.addNoRead(tz)

// 	return
// }

//func (i *indexZettelen) ReadHinweis(h hinweis.Hinweis) (mst zettel_transacted.Set, err error) {
//	if err = i.readIfNecessary(); err != nil {
//		err = errors.Wrap(err)
//		return
//	}

//	ok := false

//	if mst, ok = i.hinweisen[h]; !ok {
//		err = ErrNotFound{Id: h}
//		return
//	}

//	return
//}

//func (i *indexZettelen) ReadBezeichnung(s string) (tzs zettel_transacted.Set, err error) {
//	if err = i.readIfNecessary(); err != nil {
//		err = errors.Wrap(err)
//		return
//	}

//	ok := false

//	if tzs, ok = i.bezeichnungen[s]; !ok {
//		err = ErrNotFound{Id: stringId(s)}
//		return
//	}

//	return
//}

//func (i *indexZettelen) ReadAktenDuplicates() (tzs map[sha.Sha]zettel_transacted.Set, err error) {
//	if err = i.readIfNecessary(); err != nil {
//		err = errors.Wrap(err)
//		return
//	}

//	tzs = make(map[sha.Sha]zettel_transacted.Set, len(i.akten))

//	for s, tzsa := range i.akten {
//		if tzsa.Len() > 1 {
//			tzs[s] = tzsa
//		}
//	}

//	return
//}

//func (i *indexZettelen) ReadAkteSha(s sha.Sha) (tzs zettel_transacted.Set, err error) {
//	if err = i.readIfNecessary(); err != nil {
//		err = errors.Wrap(err)
//		return
//	}

//	ok := false

//	//TODO prevent the currently added zettels from appearing
//	if tzs, ok = i.akten[s]; !ok {
//		err = ErrNotFound{Id: s}
//		return
//	}

//	return
//}

//func (i *indexZettelen) ReadZettelSha(s sha.Sha) (tz zettel_transacted.Set, err error) {
//	if err = i.readIfNecessary(); err != nil {
//		err = errors.Wrap(err)
//		return
//	}

//	ok := false

//	if tz, ok = i.zettelen[s]; !ok {
//		err = ErrNotFound{Id: s}
//		return
//	}

//	return
//}

//func (i *indexZettelen) ReadTyp(t typ.Typ) (tzs zettel_transacted.Set, err error) {
//	if err = i.readIfNecessary(); err != nil {
//		err = errors.Wrap(err)
//		return
//	}

//	ok := false

//	if tzs, ok = i.typen[t]; !ok {
//		err = ErrNotFound{Id: t}
//		return
//	}

//	return
//}
