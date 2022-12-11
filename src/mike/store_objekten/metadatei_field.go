package store_objekten

import (
	"github.com/friedenberg/zit/src/delta/collections"
)

type LogWriter[T any] struct {
	New, Updated, Unchanged, Archived collections.WriterFunc[T]
}

//type metadateiFieldStore[
//	T gattung.Transacted[T], //transacted
//	TPtr gattung.TransactedPtr[T], //*transacted
//	O gattung.Objekte[O], //objekte
//	OPtr gattung.ObjektePtr[O], //*objekte
//	K gattung.Identifier[K], //kennung
//	KPtr gattung.IdentifierPtr[K], //*kennung
//] struct {
//	common    *common
//	logWriter LogWriter[TPtr]
//}

//func (s *metadateiFieldStore[T, TPtr, O, OPtr, K, KPtr]) SetLogWriter(
//	tlw LogWriter[TPtr],
//) {
//	s.logWriter = tlw
//}

//func makeMetadateiFieldStore[
//	T gattung.Transacted[T], //transacted
//	TPtr gattung.TransactedPtr[T], //*transacted
//	O gattung.Objekte[O], //objekte
//	OPtr gattung.ObjektePtr[O], //*objekte
//	K gattung.Identifier[K], //kennung
//	KPtr gattung.IdentifierPtr[K], //*kennung
//](
//	sa *common,
//) (s *metadateiFieldStore[T, TPtr, O, OPtr, K, KPtr], err error) {
//	s = &metadateiFieldStore[T, TPtr, O, OPtr, K, KPtr]{
//		common: sa,
//	}

//	return
//}

//func (s metadateiFieldStore[T, TPtr, O, OPtr, K, KPtr]) Flush() (err error) {
//	return
//}

//func (s metadateiFieldStore[T, TPtr, O, OPtr, K, KPtr]) transact(
//	to OPtr,
//	tk KPtr,
//) (tt T, err error) {
//	if !s.common.LockSmith.IsAcquired() {
//		err = ErrLockRequired{
//			Operation: "transact typ",
//		}

//		return
//	}

//	var mutter T

//	if mutter, err = s.ReadOne(tk); err != nil {
//		if errors.Is(err, ErrNotFound{}) {
//			err = nil
//		} else {
//			err = errors.Wrap(err)
//			return
//		}
//	}

//	//TODO-P0
//	tt = &typ.Transacted{
//		Objekte: O(*to),
//		Sku: sku.Transacted[K, KPtr]{
//			Kennung: K(*tk),
//			Schwanz: s.common.Transaktion.Time,
//		},
//	}

//	//TODO-P3 refactor into reusable
//	if mutter != nil {
//		tt.Sku.Kopf = mutter.Sku.Kopf
//		tt.Sku.Mutter[0] = mutter.Sku.Schwanz
//	} else {
//		tt.Sku.Kopf = s.common.Transaktion.Time
//	}

//	fo := objekte.MakeFormatter[T](s.common)

//	var w *age_io.Mover

//	mo := age_io.MoveOptions{
//		Age:                      s.common.Age,
//		FinalPath:                s.common.Standort.DirObjektenTypen(),
//		GenerateFinalPathFromSha: true,
//	}

//	if w, err = age_io.NewMover(mo); err != nil {
//		err = errors.Wrap(err)
//		return
//	}

//	defer errors.Deferred(&err, w.Close)

//	if _, err = fo.WriteFormat(w, tt); err != nil {
//		err = errors.Wrap(err)
//		return
//	}

//	tt.Sku.Sha = w.Sha()

//	if mutter != nil && tt.ObjekteSha().Equals(mutter.ObjekteSha()) {
//		tt = mutter

//		if err = s.logWriter.Unchanged(tt); err != nil {
//			err = errors.Wrap(err)
//			return
//		}

//		return
//	}

//	s.common.Transaktion.Add2(&tt.Sku)

//	if mutter == nil {
//		if err = s.logWriter.New(tt); err != nil {
//			err = errors.Wrap(err)
//			return
//		}
//	} else {
//		if err = s.logWriter.Updated(tt); err != nil {
//			err = errors.Wrap(err)
//			return
//		}
//	}

//	return
//}

//// TODO-P0 disambiguate from akte
//func (s metadateiFieldStore[T, TPtr, O, OPtr, K, KPtr]) WriteAkte(
//	t KPtr,
//) (err error) {
//	var w sha.WriteCloser

//	if w, err = s.common.AkteWriter(); err != nil {
//		err = errors.Wrap(err)
//		return
//	}

//	defer errors.Deferred(&err, w.Close)

//	//TODO-P0 how
//	if _, err = typ.WriteObjekteToText(w, t); err != nil {
//		err = errors.Wrap(err)
//		return
//	}

//	t.Sha = w.Sha()

//	return
//}

//func (s metadateiFieldStore[T, TPtr, O, OPtr, K, KPtr]) ReadOne(
//	k KPtr,
//) (tt T, err error) {
//	tt = s.common.Konfig.GetTyp(*k)

//	if tt == nil {
//		err = errors.Wrap(ErrNotFound{Id: k})
//		return
//	}

//	return
//}

//func (s *metadateiFieldStore[T, TPtr, O, OPtr, K, KPtr]) Create(
//	in OPtr,
//	tk KPtr,
//) (tt T, err error) {
//	if !s.common.LockSmith.IsAcquired() {
//		err = ErrLockRequired{
//			Operation: "create typ",
//		}

//		return
//	}

//	if tt, err = s.transact(in, tk); err != nil {
//		err = errors.Wrap(err)
//		return
//	}

//	return
//}

//func (s *metadateiFieldStore[T, TPtr, O, OPtr, K, KPtr]) Update(
//	t OPtr,
//	tk KPtr,
//) (tt T, err error) {
//	if !s.common.LockSmith.IsAcquired() {
//		err = ErrLockRequired{
//			Operation: "update typ",
//		}

//		return
//	}

//	if tt, err = s.transact(t, tk); err != nil {
//		err = errors.Wrap(err)
//		return
//	}

//	return
//}

//func (s metadateiFieldStore[T, TPtr, O, OPtr, K, KPtr]) AllInChain(
//	k KPtr,
//) (c []T, err error) {

//	return
//}
