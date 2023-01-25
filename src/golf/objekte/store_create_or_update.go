package objekte

//func CreateOrUpdate(
//	ls schnittstellen.LockSmith,
//	to *typ.Objekte,
//	tk *kennung.Typ,
//) (tt *typ.Transacted, err error) {
//	if !s.common.LockSmith.IsAcquired() {
//		err = ErrLockRequired{
//			Operation: "create or update typ",
//		}

//		return
//	}

//	var mutter *typ.Transacted

//	if mutter, err = s.ReadOne(tk); err != nil {
//		if errors.Is(err, ErrNotFound{}) {
//			err = nil
//		} else {
//			err = errors.Wrap(err)
//			return
//		}
//	}

//	tt = &typ.Transacted{
//		Objekte: *to,
//		Sku: sku.Transacted[kennung.Typ, *kennung.Typ]{
//			Kennung: *tk,
//			Verzeichnisse: sku.Verzeichnisse{
//				Schwanz: s.common.GetTransaktion().Time,
//			},
//		},
//	}

//	//TODO-P3 refactor into reusable
//	if mutter != nil {
//		tt.Sku.Kopf = mutter.Sku.Kopf
//		tt.Sku.Mutter[0] = mutter.Sku.Schwanz
//	} else {
//		tt.Sku.Kopf = s.common.GetTransaktion().Time
//	}

//	fo := objekte.MakeFormat[typ.Objekte, *typ.Objekte]()

//	var w *age_io.Mover

//	mo := age_io.MoveOptions{
//		Age:                      s.common.Age,
//		FinalPath:                s.common.GetStandort().DirObjektenTypen(),
//		GenerateFinalPathFromSha: true,
//	}

//	if w, err = age_io.NewMover(mo); err != nil {
//		err = errors.Wrap(err)
//		return
//	}

//	defer errors.Deferred(&err, w.Close)

//	if _, err = fo.Format(w, &tt.Objekte); err != nil {
//		err = errors.Wrap(err)
//		return
//	}

//	tt.Sku.ObjekteSha = sha.Make(w.Sha())

//	if mutter != nil && tt.GetObjekteSha().Equals(mutter.GetObjekteSha()) {
//		tt = mutter

//		if err = s.TypLogWriter.Unchanged(tt); err != nil {
//			err = errors.Wrap(err)
//			return
//		}

//		return
//	}

//	s.common.GetTransaktion().Skus.Add(&tt.Sku)
//	s.common.KonfigPtr().AddTyp(tt)

//	if mutter == nil {
//		if err = s.TypLogWriter.New(tt); err != nil {
//			err = errors.Wrap(err)
//			return
//		}
//	} else {
//		if err = s.TypLogWriter.Updated(tt); err != nil {
//			err = errors.Wrap(err)
//			return
//		}
//	}

//	return
//}
