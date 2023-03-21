package store_fs

// func (s *Store) WriteTyp(t *typ.Transacted) (te *typ.CheckedOut, err error) {
// 	te = &typ.CheckedOut{
// 		Internal: *t,
// 		External: typ.External{
// 			Sku: sku.External[kennung.Typ, *kennung.Typ]{
// 				ObjekteSha: sha.Make(t.GetObjekteSha()),
// 				Kennung:    t.Sku.Kennung,
// 				FDs: sku.ExternalFDs{
// 					Objekte: kennung.FD{
// 						Path: fmt.Sprintf("%s.%s", t.Kennung(), s.erworben.FileExtensions.Typ),
// 					},
// 				},
// 			},
// 			// TODO-P2 move to central place
// 			Objekte: t.Objekte,
// 		},
// 	}

// 	var f *os.File

// 	p := te.External.GetObjekteFD().Path

// 	if f, err = files.CreateExclusiveWriteOnly(p); err != nil {
// 		if errors.IsExist(err) {
// 			te.External, err = s.storeObjekten.Typ().ReadOneExternal(
// 				cwd.Typ{
// 					Kennung: t.Sku.Kennung,
// 					FDs: sku.ExternalFDs{
// 						Objekte: kennung.FD{
// 							Path: p,
// 						},
// 					},
// 				},
// 			)
// 		} else {
// 			err = errors.Wrap(err)
// 		}

// 		return
// 	}

// 	defer errors.Deferred(&err, f.Close)

// 	format := typ.MakeFormatText(s.storeObjekten)

// 	if _, err = format.Format(f, &te.External.Objekte); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }
