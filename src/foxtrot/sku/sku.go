package sku

// type Sku struct {
// 	WithKennung WithKennungInterface
// 	ObjekteSha  sha.Sha
// }

// func (a *Sku) SetFromSkuLike(b SkuLike) (err error) {
// 	a.WithKennung.SetMetadatei(b.GetMetadatei())
// 	a.ObjekteSha = sha.Make(b.GetObjekteSha())

// 	return
// }

// func (sk *Sku) setKennungValue(v string) (err error) {
// 	if sk.WithKennung.Kennung, err = kennung.Make(v); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

// func (sk *Sku) Set(line string) (err error) {
// 	r := strings.NewReader(line)

// 	if _, err = format.ReadSep(
// 		' ',
// 		r,
// 		ohio.MakeLineReaderIterateStrict(
// 			sk.WithKennung.Metadatei.Tai.Set,
// 			sk.WithKennung.Metadatei.Gattung.Set,
// 			sk.setKennungValue,
// 			sk.ObjekteSha.Set,
// 			sk.WithKennung.Metadatei.AkteSha.Set,
// 		),
// 	); err != nil {
// 		if err1 := sk.setOld(line); err1 != nil {
// 			err = errors.MakeMulti(err, err1)
// 			return
// 		}

// 		err = nil

// 		return
// 	}

// 	return
// }

// func (sk *Sku) setOld(line string) (err error) {
// 	r := strings.NewReader(line)

// 	if _, err = format.ReadSep(
// 		' ',
// 		r,
// 		ohio.MakeLineReaderIterateStrict(
// 			sk.WithKennung.Metadatei.Gattung.Set,
// 			sk.WithKennung.Metadatei.Tai.Set,
// 			sk.setKennungValue,
// 			sk.ObjekteSha.Set,
// 			sk.WithKennung.Metadatei.AkteSha.Set,
// 		),
// 	); err != nil {
// 		err = errors.Wrapf(err, "Sku2: %s", line)
// 		return
// 	}

// 	return
// }

// func (a *Sku) ResetWith(b Sku) {
// 	errors.TodoP4("should these be more ResetWith calls?")
// 	a.WithKennung.Metadatei.Gattung = b.WithKennung.Metadatei.Gattung
// 	a.WithKennung.Metadatei.Tai = b.WithKennung.Metadatei.Tai
// 	a.WithKennung.Kennung = b.WithKennung.Kennung
// 	a.ObjekteSha = b.ObjekteSha
// 	a.WithKennung.Metadatei.AkteSha = b.WithKennung.Metadatei.AkteSha
// }

// func (a *Sku) Reset() {
// 	a.WithKennung.Metadatei.Gattung.Reset()
// 	a.WithKennung.Metadatei.Tai.Reset()

// 	kp := a.WithKennung.Kennung.KennungPtrClone()
// 	kp.Reset()
// 	a.WithKennung.Kennung = kp.KennungClone()

// 	a.ObjekteSha.Reset()
// 	a.WithKennung.Metadatei.AkteSha.Reset()
// }

// func (a Sku) GetMetadatei() Metadatei {
// 	return a.WithKennung.Metadatei
// }

// func (a *Sku) GetMetadateiPtr() *Metadatei {
// 	return &a.WithKennung.Metadatei
// }

// func (a Sku) GetTai() kennung.Tai {
// 	return a.WithKennung.Metadatei.Tai
// }

// func (a Sku) GetKey() string {
// 	return a.String()
// }

// func (a Sku) GetTime() kennung.Time {
// 	return a.WithKennung.Metadatei.Tai.AsTime()
// }

// func (a Sku) GetId() Kennung {
// 	return a.WithKennung.Kennung
// }

// func (a Sku) GetKennung() kennung.Kennung {
// 	return a.WithKennung.Kennung
// }

// func (a Sku) GetGattung() schnittstellen.GattungLike {
// 	return a.WithKennung.Metadatei.Gattung
// }

// func (a Sku) GetObjekteSha() schnittstellen.ShaLike {
// 	return a.ObjekteSha
// }

// func (a Sku) GetAkteSha() schnittstellen.ShaLike {
// 	return a.WithKennung.Metadatei.AkteSha
// }

// func (a Sku) Less(b Sku) (ok bool) {
// 	if a.WithKennung.Metadatei.Tai.Less(b.WithKennung.Metadatei.Tai) {
// 		ok = true
// 		return
// 	}

// 	return
// }

// func (a Sku) EqualsSkuLike(b SkuLike) (ok bool) {
// 	return values.Equals(a, b)
// }

// func (a Sku) EqualsAny(b any) (ok bool) {
// 	return values.Equals(a, b)
// }

// func (a Sku) Equals(b Sku) (ok bool) {
// 	if a != b {
// 		return false
// 	}

// 	return true
// }

// func (s Sku) String() string {
// 	return fmt.Sprintf(
// 		"%s %s %s %s %s",
// 		s.WithKennung.Metadatei.Tai,
// 		s.WithKennung.Metadatei.Gattung,
// 		s.WithKennung.Kennung,
// 		s.ObjekteSha,
// 		s.WithKennung.Metadatei.AkteSha,
// 	)
// }
