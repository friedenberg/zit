package sku

// func init() {
// 	gob.Register(&Transacted[kennung.Hinweis, *kennung.Hinweis]{})
// 	gob.Register(&Transacted[kennung.Etikett, *kennung.Etikett]{})
// 	gob.Register(&Transacted[kennung.Typ, *kennung.Typ]{})
// 	gob.Register(&Transacted[kennung.Kasten, *kennung.Kasten]{})
// 	gob.Register(&Transacted[kennung.Konfig, *kennung.Konfig]{})
// }

// // TODO-P2 move sku.Sku to sku.Transacted
// type Transacted[K kennung.KennungLike[K], KPtr kennung.KennungLikePtr[K]]
// struct {
// 	Kennung          K
// 	Metadatei        metadatei.Metadatei
// 	ObjekteSha       sha.Sha
// 	TransactionIndex values.Int
// 	Kopf             kennung.Tai
// }

// func (t *Transacted[K, KPtr]) SetFromSkuLike(sk SkuLike) (err error) {
// 	if err = KPtr(&t.Kennung).Set(sk.GetKennungLike().String()); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	t.ObjekteSha = sha.Make(sk.GetObjekteSha())
// 	t.Metadatei.ResetWith(sk.GetMetadatei())
// 	t.GetMetadateiPtr().Tai = sk.GetTai()

// 	t.Kopf = sk.GetTai()

// 	return
// }

// func MakeSkuLikeSansObjekteSha(
// 	m metadatei.Metadatei,
// 	k kennung.Kennung,
// ) (sk SkuLikePtr, err error) {
// 	switch kt := k.(type) {
// 	case *kennung.Hinweis:
// 		sk = &Transacted[kennung.Hinweis, *kennung.Hinweis]{
// 			Metadatei: m,
// 			Kennung:   *kt,
// 		}

// 	case *kennung.Etikett:
// 		sk = &Transacted[kennung.Etikett, *kennung.Etikett]{
// 			Metadatei: m,
// 			Kennung:   *kt,
// 		}

// 	case *kennung.Typ:
// 		sk = &Transacted[kennung.Typ, *kennung.Typ]{
// 			Metadatei: m,
// 			Kennung:   *kt,
// 		}

// 	case *kennung.Kasten:
// 		sk = &Transacted[kennung.Kasten, *kennung.Kasten]{
// 			Metadatei: m,
// 			Kennung:   *kt,
// 		}

// 	case *kennung.Konfig:
// 		sk = &Transacted[kennung.Konfig, *kennung.Konfig]{
// 			Metadatei: m,
// 			Kennung:   *kt,
// 		}

// 	default:
// 		err = errors.Errorf("unsupported kennung: %T -> %q", kt, kt)
// 		return
// 	}

// 	return
// }

// func MakeSkuLike(
// 	m metadatei.Metadatei,
// 	k kennung.Kennung,
// 	os sha.Sha,
// ) (sk *Transacted2, err error) {
// 	if sk, err = MakeSkuLikeSansObjekteSha2(m, k); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	sk.SetObjekteSha(os)

// 	return
// }

// func (a Transacted[K, KPtr]) ImmutableClone() SkuLike {
// 	return a
// }

// func (a Transacted[K, KPtr]) MutableClone() SkuLikePtr {
// 	return &a
// }

// func (a Transacted[K, KPtr]) String() string {
// 	return fmt.Sprintf(
// 		"%s %s %s",
// 		a.Kennung,
// 		a.ObjekteSha,
// 		a.Metadatei.AkteSha,
// 	)
// }

// func (a Transacted[K, KPtr]) GetSkuLike() SkuLike {
// 	return a
// }

// func (a *Transacted[K, KPtr]) GetSkuLikePtr() SkuLikePtr {
// 	return a
// }

// func (a Transacted[K, KPtr]) GetEtiketten() kennung.EtikettSet {
// 	return a.Metadatei.GetEtiketten()
// }

// func (a Transacted[K, KPtr]) GetTyp() kennung.Typ {
// 	return a.Metadatei.Typ
// }

// func (a Transacted[K, KPtr]) GetMetadatei() metadatei.Metadatei {
// 	return a.Metadatei
// }

// func (a *Transacted[K, KPtr]) GetMetadateiPtr() *metadatei.Metadatei {
// 	return &a.Metadatei
// }

// func (a *Transacted[K, KPtr]) SetMetadatei(m metadatei.Metadatei) {
// 	a.Metadatei = m
// }

// func (a Transacted[K, KPtr]) GetTai() kennung.Tai {
// 	return a.GetMetadatei().GetTai()
// }

// func (a Transacted[K, KPtr]) GetKopf() kennung.Tai {
// 	return a.Kopf
// }

// func (a *Transacted[K, KPtr]) SetTai(t kennung.Tai) {
// 	a.GetMetadateiPtr().Tai = t
// }

// func (a Transacted[K, KPtr]) GetKennung() K {
// 	return a.Kennung
// }

// func (a *Transacted[K, KPtr]) GetKennungPtr() KPtr {
// 	return &a.Kennung
// }

// func (a Transacted[K, KPtr]) GetKennungLike() kennung.Kennung {
// 	return a.Kennung
// }

// func (a *Transacted[K, KPtr]) GetKennungLikePtr() kennung.KennungPtr {
// 	return KPtr(&a.Kennung)
// }

// func (a *Transacted[K, KPtr]) SetKennungLike(kl kennung.Kennung) (err error)
// {
// 	switch k := kl.(type) {
// 	case K:
// 		a.Kennung = k

// 	case KPtr:
// 		a.Kennung = K(*k)

// 	default:
// 		err = errors.Errorf("expected kennung of type %T but got %T: %q",
// a.Kennung, k, kl)
// 		return
// 	}

// 	return
// }

// func (a Transacted[K, KPtr]) GetExternal() External[K, KPtr] {
// 	return External[K, KPtr]{
// 		Transacted: a,
// 	}
// }

// func (a *Transacted[K, KPtr]) SetTransactionIndex(i int) {
// 	a.TransactionIndex.SetInt(i)
// }

// func (a *Transacted[K, KPtr]) Reset() {
// 	a.Kopf.Reset()
// 	a.ObjekteSha.Reset()
// 	KPtr(&a.Kennung).Reset()
// 	a.Metadatei.Reset()
// 	a.TransactionIndex.Reset()
// }

// func (a *Transacted[K, KPtr]) ResetWith(b Transacted[K, KPtr]) {
// 	a.Kopf = b.Kopf
// 	a.ObjekteSha = b.ObjekteSha
// 	a.Kennung = b.Kennung
// 	a.Metadatei.ResetWith(b.Metadatei)
// 	a.TransactionIndex.SetInt(b.TransactionIndex.Int())
// }

// // TODO-P2 switch this to default
// func (a *Transacted[T2, T3]) ResetWithPtr(
// 	b *Transacted[T2, T3],
// ) {
// 	a.ResetWith(*b)
// }

// func (a Transacted[K, KPtr]) Less(b Transacted[K, KPtr]) (ok bool) {
// 	if a.GetTai().Less(b.GetTai()) {
// 		ok = true
// 		return
// 	}

// 	// if a.GetTai().Equals(b.GetTai()) &&
// 	// 	a.TransactionIndex.Less(b.TransactionIndex) {
// 	// 	ok = true
// 	// 	return
// 	// }

// 	return
// }

// func (a Transacted[K, KPtr]) EqualsSkuLike(b SkuLike) bool {
// 	return values.Equals(a, b) || values.EqualsPtr(a, b)
// }

// func (a Transacted[K, KPtr]) EqualsAny(b any) (ok bool) {
// 	return values.Equals(a, b)
// }

// func (a Transacted[K, KPtr]) Equals(b Transacted[K, KPtr]) (ok bool) {
// 	if !a.TransactionIndex.Equals(b.TransactionIndex) {
// 		return
// 	}

// 	if a.GetKennung().String() != b.GetKennung().String() {
// 		return
// 	}

// 	// TODO-P2 determine why objekte shas in import test differed
// 	// if !a.ObjekteSha.Equals(b.ObjekteSha) {
// 	// 	return
// 	// }

// 	if !a.Metadatei.Equals(b.Metadatei) {
// 		return
// 	}

// 	return true
// }

// func (s Transacted[K, KPtr]) GetGattung() schnittstellen.GattungLike {
// 	return s.Kennung.GetGattung()
// }

// func (s *Transacted[K, KPtr]) IsNew() bool {
// 	return s.Metadatei.Verzeichnisse.Mutter.IsNull()
// }

// func (s *Transacted[K, KPtr]) SetObjekteSha(v schnittstellen.ShaLike) {
// 	s.ObjekteSha = sha.Make(v)
// }

// func (s Transacted[K, KPtr]) GetObjekteSha() schnittstellen.ShaLike {
// 	return s.ObjekteSha
// }

// func (s Transacted[K, KPtr]) GetAkteSha() schnittstellen.ShaLike {
// 	return s.Metadatei.AkteSha
// }

// func (s *Transacted[K, KPtr]) SetAkteSha(sh schnittstellen.ShaLike) {
// 	s.Metadatei.AkteSha = sha.Make(sh)
// }

// func (s Transacted[K, KPtr]) GetTransactionIndex() values.Int {
// 	return s.TransactionIndex
// }

// func (o Transacted[K, KPtr]) GetKey() string {
// 	return kennung.FormattedString(o.GetKennung())
// }
