package objekte_store

// type Delegate[T any] struct {
// 	New       schnittstellen.FuncIter[T]
// 	Updated   schnittstellen.FuncIter[T]
// 	Unchanged schnittstellen.FuncIter[T]
// }

// type reindexer[
// 	T schnittstellen.Objekte[T],
// 	T1 schnittstellen.ObjektePtr[T],
// 	T2 schnittstellen.Id[T2],
// 	T3 schnittstellen.IdPtr[T2],
// 	T4 any,
// 	T5 schnittstellen.VerzeichnissePtr[T4, T],
// ] struct {
// 	clock    ts.Clock
// 	ls       schnittstellen.LockSmith
// 	oaf      schnittstellen.ObjekteAkteWriterFactory
// 	reader   TransactedReader[T3, *objekte.Transacted[T, T1, T2, T3, T4, T5]]
// 	delegate Delegate[*objekte.Transacted[T, T1, T2, T3, T4, T5]]
// }

// func MakeReindexer[
// 	T schnittstellen.Objekte[T],
// 	T1 schnittstellen.ObjektePtr[T],
// 	T2 schnittstellen.Id[T2],
// 	T3 schnittstellen.IdPtr[T2],
// 	T4 any,
// 	T5 schnittstellen.VerzeichnissePtr[T4, T],
// ](
// 	clock ts.Clock,
// 	ls schnittstellen.LockSmith,
// 	oaf schnittstellen.ObjekteAkteWriterFactory,
// 	reader TransactedReader[T3, *objekte.Transacted[T, T1, T2, T3, T4, T5]],
// 	delegate Delegate[*objekte.Transacted[T, T1, T2, T3, T4, T5]],
// ) (cou *reindexer[T, T1, T2, T3, T4, T5]) {
// 	return &reindexer[T, T1, T2, T3, T4, T5]{
// 		clock:    clock,
// 		ls:       ls,
// 		oaf:      oaf,
// 		reader:   reader,
// 		delegate: delegate,
// 	}
// }

// func (r reindexer[T, T1, T2, T3, T4, T5]) ReindexOne(
// 	sk sku.DataIdentity,
// ) (o schnittstellen.Stored, err error) {
// 	var t *objekte.Transacted[T, T1, T2, T3, T4, T5]

// 	if t, err = s.InflateFromDataIdentity(sk); err != nil {
// 		if errors.Is(err, toml.Error{}) {
// 			err = nil
// 			return
// 		} else {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	}

// 	o = t

// 	if te.IsNew() {
// 		if err = r.delegate.New(te); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	} else {
// 		if err = r.delegate.Updated(te); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	}

// 	return
// }
