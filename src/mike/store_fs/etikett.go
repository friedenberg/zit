package store_fs

// func (s *Store) WriteEtikett(t *etikett.Transacted) (te *etikett.CheckedOut, err error) {
// 	te = &etikett.CheckedOut{
// 		Internal: *t,
// 		External: etikett.External{
// 			Sku: sku.External[kennung.Etikett, *kennung.Etikett]{
// 				ObjekteSha: sha.Make(t.GetObjekteSha()),
// 				Kennung:    t.Sku.Kennung,
// 				FDs: sku.ExternalFDs{
// 					Objekte: kennung.FD{
// 						Path: fmt.Sprintf("%s.%s", t.Kennung(), s.erworben.FileExtensions.Etikett),
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
// 			te.External, err = s.storeObjekten.Etikett().ReadOneExternal(
// 				cwd.Etikett{
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

// 	format := etikett.MakeFormatText(s.storeObjekten)

// 	if _, err = format.Format(f, &te.External.Objekte); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }
