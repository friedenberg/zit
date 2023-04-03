package store_objekten

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/hinweisen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/kilo/store_util"
	"github.com/friedenberg/zit/src/kilo/store_verzeichnisse"
)

type ZettelStore interface {
	CommonStoreBase[
		zettel.Objekte,
		*zettel.Objekte,
		kennung.Hinweis,
		*kennung.Hinweis,
		zettel.Verzeichnisse,
		*zettel.Verzeichnisse,
	]

	objekte_store.Creator[
		zettel.Objekte,
		*zettel.Transacted,
	]

	objekte_store.CheckedOutUpdater[
		zettel.CheckedOut,
		*zettel.Transacted,
	]

	objekte_store.Updater[
		*zettel.Objekte,
		*kennung.Hinweis,
		// schnittstellen.Value,
		*zettel.Transacted,
	]

	WriteZettelObjekte(z zettel.Objekte) (sh sha.Sha, err error)
}

type zettelStore struct {
	*commonStore[
		zettel.Objekte,
		*zettel.Objekte,
		kennung.Hinweis,
		*kennung.Hinweis,
		zettel.Verzeichnisse,
		*zettel.Verzeichnisse,
	]

	format      zettel.ObjekteFormat
	protoZettel zettel.ProtoZettel

	verzeichnisseSchwanzen *verzeichnisseSchwanzen
	verzeichnisseAll       *store_verzeichnisse.Zettelen
}

func makeZettelStore(
	sa store_util.StoreUtil,
	p schnittstellen.Pool[zettel.Transacted, *zettel.Transacted],
) (s *zettelStore, err error) {
	s = &zettelStore{
		protoZettel: zettel.MakeProtoZettel(sa.GetKonfig()),
		format: zettel.MakeObjekteTextFormat(
			sa,
			nil,
		),
	}

	s.commonStore, err = makeCommonStore[
		zettel.Objekte,
		*zettel.Objekte,
		kennung.Hinweis,
		*kennung.Hinweis,
		zettel.Verzeichnisse,
		*zettel.Verzeichnisse,
	](
		s,
		sa,
		s,
		&zettel.FormatObjekte{
			IgnoreTypErrors: true,
		},
		nil,
		nil,
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.verzeichnisseSchwanzen, err = makeVerzeichnisseSchwanzen(
		s.StoreUtil,
		p,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.verzeichnisseAll, err = store_verzeichnisse.MakeZettelen(
		s.StoreUtil.GetKonfig(),
		s.StoreUtil.GetStandort().DirVerzeichnisseZettelenNeue(),
		s.StoreUtil,
		p,
		nil,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) Flush() (err error) {
	if err = s.verzeichnisseSchwanzen.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.verzeichnisseAll.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s zettelStore) WriteZettelObjekte(z zettel.Objekte) (sh sha.Sha, err error) {
	// no lock required

	var wc sha.WriteCloser

	if wc, err = s.StoreUtil.ObjekteWriter(gattung.Zettel); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, wc)

	c := zettel.ObjekteFormatterContext{
		Zettel: z,
	}

	f := zettel.FormatObjekte{}

	if _, err = f.Format(wc, &c.Zettel); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = sha.Make(wc.Sha())

	return
}

func (s *zettelStore) writeNamedZettelToIndex(
	tz *zettel.Transacted,
) (err error) {
	errors.Log().Print("writing to index")

	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "write named zettel to index",
		}

		return
	}

	errors.Log().Printf("writing zettel to index: %s", tz.Sku)

	if err = s.verzeichnisseSchwanzen.Add(tz, tz.Sku.Kennung.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.verzeichnisseAll.Add(tz, tz.Sku.ObjekteSha.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.GetKennungIndex().AddHinweis(tz.Sku.Kennung); err != nil {
		if errors.Is(err, hinweisen.ErrDoesNotExist{}) {
			errors.Log().Printf("kennung does not contain value: %s", err)
			err = nil
		} else {
			err = errors.Wrapf(err, "failed to write zettel to index: %s", tz.Sku)
			return
		}
	}

	return
}

func (s *zettelStore) ReadOneExternal(
	e cwd.Zettel,
) (ez zettel.External, err error) {
	// support akte
	todo.Implement()

	ez.Sku.ResetWithExternalMaybe(e)

	c := zettel.ObjekteParserContext{}

	var f *os.File

	if f, err = files.Open(ez.GetObjekteFD().Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	if _, err = s.format.Parse(f, &c); err != nil {
		err = errors.Wrapf(err, "%s", f.Name())
		return
	}

	if ez.Sku.ObjekteSha, err = s.WriteZettelObjekte(
		c.Zettel,
	); err != nil {
		err = errors.Wrapf(err, "%s", f.Name())
		return
	}

	ez.Objekte = c.Zettel
	ez.Sku.FDs.Akte.Path = c.AktePath

	unrecoverableErrors := errors.MakeMulti()

	for _, e := range errors.Split(c.Errors) {
		var err1 zettel.ErrHasInvalidAkteShaOrFilePath

		if errors.As(e, &err1) {
			var mutter *zettel.Transacted

			if mutter, err = s.ReadOne(
				&ez.Sku.Kennung,
			); err != nil {
				unrecoverableErrors.Add(errors.Wrap(err))
				continue
			}

			ez.Objekte.Akte = mutter.Objekte.Akte

			continue
		}

		unrecoverableErrors.Add(e)
	}

	if !unrecoverableErrors.Empty() {
		err = errors.Wrap(unrecoverableErrors)
		return
	}

	return
}

func (s zettelStore) ReadOne(
	i *kennung.Hinweis,
) (tz *zettel.Transacted, err error) {
	if tz, err = s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(*i); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *zettelStore) ReadAllSchwanzen(
	w schnittstellen.FuncIter[*zettel.Transacted],
) (err error) {
	return i.verzeichnisseSchwanzen.ReadMany(w)
}

func (i *zettelStore) ReadAll(
	w schnittstellen.FuncIter[*zettel.Transacted],
) (err error) {
	return i.verzeichnisseAll.ReadMany(w)
}

func (s *zettelStore) Create(
	in zettel.Objekte,
) (tz *zettel.Transacted, err error) {
	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "create",
		}

		return
	}

	// if in.IsEmpty() || s.protoZettel.Equals(in) {
	// 	err = errors.Normalf("zettel is empty")
	// 	return
	// }

	s.protoZettel.Apply(&in)

	if err = s.StoreUtil.GetKonfig().ApplyToMetadatei(&in.Metadatei, in.Typ); err != nil {
		err = errors.Wrap(err)
		return
	}

	var shaObj sha.Sha

	if shaObj, err = s.WriteZettelObjekte(in); err != nil {
		err = errors.Wrap(err)
		return
	}

	// If the zettel exists, short circuit and return that
	todo.Implement()
	// if tz2, err2 := s.ReadOne(shaObj); err2 == nil {
	// 	tz = tz2
	// 	return
	// }

	var ken kennung.Hinweis

	if ken, err = s.StoreUtil.GetKennungIndex().CreateHinweis(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if tz, err = s.addZettelToTransaktion(
		&in,
		&shaObj,
		&ken,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.GetKennungIndex().Add(tz.Objekte.Metadatei.Etiketten); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.AddMatchable(tz); err != nil {
		err = errors.Wrapf(err, "failed to write zettel to index: %s", tz.Sku)
		return
	}

	errors.TodoP2("assert no changes")
	if err = s.LogWriter.New(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) updateExternal(
	co objekte.ExternalLike,
) (tl objekte.TransactedLike, err error) {
	ze := co.(*zettel.External)
	return s.Update(&ze.Objekte, &ze.Sku.Kennung)
}

func (s *zettelStore) UpdateCheckedOut(
	co zettel.CheckedOut,
) (t *zettel.Transacted, err error) {
	errors.TodoP2("support dry run")

	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "update",
		}

		return
	}

	if err = s.StoreUtil.GetKonfig().ApplyToMetadatei(
		&co.External.Objekte.Metadatei,
		co.External.Objekte.Typ,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var shaObj sha.Sha

	if shaObj, err = s.WriteZettelObjekte(co.External.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	if shaObj.Equals(co.Internal.Sku.ObjekteSha) {
		t = &co.Internal

		if err = s.LogWriter.Unchanged(t); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if t, err = s.addZettelToTransaktion(
		&co.External.Objekte,
		&shaObj,
		&co.External.Sku.Kennung,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.writeNamedZettelToIndex(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.AddMatchable(t); err != nil {
		err = errors.Wrapf(err, "failed to write zettel to index: %s", t.Sku)
		return
	}

	if err = s.LogWriter.Updated(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) Update(
	z *zettel.Objekte,
	h *kennung.Hinweis,
) (tz *zettel.Transacted, err error) {
	errors.TodoP2("support dry run")

	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "update",
		}

		return
	}

	if err = s.StoreUtil.GetKonfig().ApplyToMetadatei(
		&z.Metadatei,
		z.Typ,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var mutter *zettel.Transacted

	if mutter, err = s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(
		*h,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var shaObj sha.Sha

	if shaObj, err = s.WriteZettelObjekte(*z); err != nil {
		err = errors.Wrap(err)
		return
	}

	if shaObj.Equals(mutter.Sku.ObjekteSha) {
		tz = mutter

		if err = s.LogWriter.Unchanged(tz); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if tz, err = s.addZettelToTransaktion(
		z,
		&shaObj,
		h,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.AddMatchable(tz); err != nil {
		err = errors.Wrapf(err, "failed to write zettel to index: %s", tz.Sku)
		return
	}

	if err = s.LogWriter.Updated(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// func (s zettelStore) AllInChain(h hinweis.Hinweis) (c []*zettel.Transacted, err error) {
// 	mst := zettel.MakeMutableSetUnique(0)

// 	if err = s.verzeichnisseAll.ReadMany(
// 		func(z *zettel.Transacted) (err error) {
// 			if !z.Sku.Kennung.Equals(h) {
// 				err = collections.ErrStopIteration
// 				return
// 			}

// 			return
// 		},
// 		mst.AddAndDoNotRepool,
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	c = mst.Elements()

// 	sort.Slice(
// 		c,
// 		func(i, j int) bool { return c[i].Sku.Less(&c[j].Sku) },
// 	)

// 	return
// }

func (s *zettelStore) addZettelToTransaktion(
	zo *zettel.Objekte,
	zs *sha.Sha,
	zk *kennung.Hinweis,
) (tz *zettel.Transacted, err error) {
	errors.Log().Printf("adding zettel to transaktion: %s", zk)

	if tz, err = s.transactedWithHead(
		*zo,
		*zk,
		s.StoreUtil.GetTransaktionStore().GetTransaktion(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	tz.Sku.Kennung = *zk
	tz.Sku.ObjekteSha = *zs

	s.StoreUtil.CommitTransacted(tz)

	return
}

// TODO-P1 modify this to not require immediate mutter
// should only be called when moving forward through time, as there is a
// dependency on the index being accurate for the immediate mutter of the zettel
// in the arguments
func (s *zettelStore) transactedWithHead(
	z zettel.Objekte,
	h kennung.Hinweis,
	t *transaktion.Transaktion,
) (tz *zettel.Transacted, err error) {
	tz = &zettel.Transacted{
		Objekte: z,
		Sku: sku.Transacted[kennung.Hinweis, *kennung.Hinweis]{
			Kennung: h,
			Verzeichnisse: sku.Verzeichnisse{
				Kopf:    t.Time,
				Schwanz: t.Time,
			},
		},
	}

	tz.Verzeichnisse.ResetWithObjekte(tz.Objekte)

	return
}

func (s *zettelStore) Inherit(tz *zettel.Transacted) (err error) {
	errors.Log().Printf("inheriting %s", tz.Sku.ObjekteSha)

	if _, err = s.WriteZettelObjekte(tz.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.StoreUtil.CommitTransacted(tz)

	errExists := s.StoreUtil.GetAbbrStore().HinweisExists(tz.Sku.Kennung)

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.LogWriter.NewOrUpdated(errExists)(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) ReindexOne(
	sk sku.DataIdentity,
) (o kennung.Matchable, err error) {
	var tz *zettel.Transacted
	defer s.pool.Put(tz)

	errors.Log().Printf("reindexing: %#v", o)

	if tz, err = s.InflateFromDataIdentity(sk); err != nil {
		// TODO-P2 decide on how to handle format errors
		errors.Err().Print(err)
		err = nil
		// err = errors.Wrap(err)
		return
	}

	o = tz
	errExists := s.StoreUtil.GetAbbrStore().HinweisExists(tz.Sku.Kennung)

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.LogWriter.NewOrUpdated(errExists)(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
