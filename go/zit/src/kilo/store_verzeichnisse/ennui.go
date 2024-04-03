package store_verzeichnisse

// func (i *Store) GetEnnuiShas() ennui.Ennui {
// 	return i.ennuiShas
// }

// func (i *Store) GetEnnuiKennung() ennui.Ennui {
// 	return i.ennuiKennung
// }

// func (i *Store) ExistsOneSha(sh *sha.Sha) (err error) {
// 	if _, err = i.ennuiShas.ReadOne(sh); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	err = collections.ErrExists

// 	return
// }

// func (i *Store) ReadOneShas(sh *sha.Sha) (out *sku.Transacted, err error) {
// 	var loc ennui.Loc

// 	if loc, err = i.ennuiShas.ReadOne(sh); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return i.readLoc(loc)
// }

// func (i *Store) ReadOneKennung(
// 	h schnittstellen.Stringer,
// ) (out *sku.Transacted, err error) {
// 	sh := sha.FromString(h.String())
// 	defer sha.GetPool().Put(sh)

// 	var loc ennui.Loc

// 	if loc, err = i.ennuiKennung.ReadOne(sh); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return i.readLoc(loc)
// }

// func (i *Store) ReadOneAll(
// 	mg metadatei.Getter,
// 	kennungPtr kennung.Kennung,
// ) (out []ennui.Loc, err error) {
// 	var locKennung ennui.Loc

// 	wg := iter.MakeErrorWaitGroupParallel()

// 	wg.Do(func() (err error) {
// 		sh := sha.FromString(kennungPtr.String())
// 		defer sha.GetPool().Put(sh)

// 		if locKennung, err = i.ennuiKennung.ReadOne(sh); err != nil {
// 			if collections.IsErrNotFound(err) {
// 				err = nil
// 			} else {
// 				err = errors.Wrap(err)
// 			}

// 			return
// 		}

// 		return
// 	})

// 	wg.Do(func() (err error) {
// 		if err = i.ennuiShas.ReadAll(mg.GetMetadatei(), &out); err != nil {
// 			if collections.IsErrNotFound(err) {
// 				err = nil
// 			} else {
// 				err = errors.Wrap(err)
// 			}

// 			return
// 		}

// 		return
// 	})

// 	if err = wg.GetError(); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if !locKennung.IsEmpty() {
// 		out = append(out, locKennung)
// 	}

// 	return
// }

// func (i *Store) readLoc(loc ennui.Loc) (sk *sku.Transacted, err error) {
// 	p := &i.pages[loc.Page]

// 	var f *os.File

// 	if f, err = files.OpenFile(
// 		p.Path(),
// 		os.O_RDONLY,
// 		0o666,
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	defer errors.DeferredCloser(&err, f)

// 	coder := binaryDecoder{
// 		QueryGroup: &sigil{Sigil: kennung.SigilAll},
// 	}

// 	sk = sku.GetTransactedPool().Get()

// 	if _, err = coder.readFormatExactly(
// 		f,
// 		loc,
// 		&Sku{skuWithSigil: skuWithSigil{Transacted: sk}},
// 	); err != nil {
// 		sku.GetTransactedPool().Put(sk)
// 		sk = nil
// 		err = errors.Wrapf(err, "%s", loc)
// 		return
// 	}

// 	if sk != nil {
// 		log.Debug().Print(sk.Metadatei.GetEtiketten())
// 	}

// 	return
// }
