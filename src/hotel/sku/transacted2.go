package sku

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type Transacted2 struct {
	Kennung          kennung.Kennung2
	Metadatei        metadatei.Metadatei
	ObjekteSha       sha.Sha
	TransactionIndex values.Int
	Kopf             kennung.Tai
}

func (t *Transacted2) SetFromSkuLike(sk SkuLike) (err error) {
	err = t.Kennung.SetWithGattung(
		sk.GetKennungLike().String(),
		sk.GetGattung(),
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	t.ObjekteSha = sha.Make(sk.GetObjekteSha())
	t.Metadatei.ResetWith(sk.GetMetadatei())
	t.GetMetadateiPtr().Tai = sk.GetTai()

	t.Kopf = sk.GetTai()

	return
}

func MakeSkuLikeSansObjekteSha2(
	m metadatei.Metadatei,
	k kennung.Kennung,
) (sk *Transacted2, err error) {
	sk = &Transacted2{
		Metadatei: m,
	}

	if err = sk.Kennung.SetWithKennung(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeSkuLike2(
	m metadatei.Metadatei,
	k kennung.KennungPtr,
	os sha.Sha,
) (sk SkuLikePtr, err error) {
	if sk, err = MakeSkuLikeSansObjekteSha2(m, k); err != nil {
		err = errors.Wrap(err)
		return
	}

	sk.SetObjekteSha(os)

	return
}

func (a Transacted2) ImmutableClone() SkuLike {
	return a
}

func (a Transacted2) MutableClone() SkuLikePtr {
	return &a
}

func (a Transacted2) String() string {
	return fmt.Sprintf(
		"%s %s %s",
		a.Kennung,
		a.ObjekteSha,
		a.Metadatei.AkteSha,
	)
}

func (a Transacted2) GetSkuLike() SkuLike {
	return a
}

func (a *Transacted2) GetSkuLikePtr() SkuLikePtr {
	return a
}

func (a Transacted2) GetEtiketten() kennung.EtikettSet {
	return a.Metadatei.GetEtiketten()
}

func (a Transacted2) GetTyp() kennung.Typ {
	return a.Metadatei.Typ
}

func (a Transacted2) GetMetadatei() metadatei.Metadatei {
	return a.Metadatei
}

func (a *Transacted2) GetMetadateiPtr() *metadatei.Metadatei {
	return &a.Metadatei
}

func (a *Transacted2) SetMetadatei(m metadatei.Metadatei) {
	a.Metadatei = m
}

func (a Transacted2) GetTai() kennung.Tai {
	return a.GetMetadatei().GetTai()
}

func (a Transacted2) GetKopf() kennung.Tai {
	return a.Kopf
}

func (a *Transacted2) SetTai(t kennung.Tai) {
	a.GetMetadateiPtr().Tai = t
}

func (a Transacted2) GetKennung() kennung.Kennung {
	return a.Kennung
}

func (a *Transacted2) GetKennungPtr() kennung.KennungPtr {
	return &a.Kennung
}

func (a Transacted2) GetKennungLike() kennung.Kennung {
	return a.Kennung
}

func (a *Transacted2) GetKennungLikePtr() kennung.KennungPtr {
	return &a.Kennung
}

func (a *Transacted2) SetKennungLike(kl kennung.Kennung) (err error) {
	if err = a.Kennung.Set(kl.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Transacted2) Reset() {
	a.Kopf.Reset()
	a.ObjekteSha.Reset()

	// TODO-P2 remove in favor of kennung pkg
	if a.Kennung.KennungPtr != nil {
		a.Kennung.Reset()
	}

	a.Metadatei.Reset()
	a.TransactionIndex.Reset()
}

func (a *Transacted2) ResetWith(b Transacted2) {
	a.Kopf = b.Kopf
	a.ObjekteSha = b.ObjekteSha
	a.Kennung.ResetWithKennung(b.Kennung)
	a.Metadatei.ResetWith(b.Metadatei)
	a.TransactionIndex.SetInt(b.TransactionIndex.Int())
}

// TODO-P2 switch this to default
func (a *Transacted2) ResetWithPtr(
	b *Transacted2,
) {
	a.ResetWith(*b)
}

func (a Transacted2) Less(b Transacted2) (ok bool) {
	if a.GetTai().Less(b.GetTai()) {
		ok = true
		return
	}

	// if a.GetTai().Equals(b.GetTai()) &&
	// 	a.TransactionIndex.Less(b.TransactionIndex) {
	// 	ok = true
	// 	return
	// }

	return
}

func (a Transacted2) EqualsSkuLike(b SkuLike) bool {
	return values.Equals(a, b) || values.EqualsPtr(a, b)
}

func (a Transacted2) EqualsAny(b any) (ok bool) {
	return values.Equals(a, b)
}

func (a Transacted2) Equals(b Transacted2) (ok bool) {
	if !a.TransactionIndex.Equals(b.TransactionIndex) {
		return
	}

	if a.GetKennung().String() != b.GetKennung().String() {
		return
	}

	// TODO-P2 determine why objekte shas in import test differed
	// if !a.ObjekteSha.Equals(b.ObjekteSha) {
	// 	return
	// }

	if !a.Metadatei.Equals(b.Metadatei) {
		return
	}

	return true
}

func (s Transacted2) GetGattung() schnittstellen.GattungLike {
	return s.Kennung.GetGattung()
}

func (s *Transacted2) IsNew() bool {
	return s.Metadatei.Verzeichnisse.Mutter.IsNull()
}

func (s *Transacted2) SetObjekteSha(v schnittstellen.ShaLike) {
	s.ObjekteSha = sha.Make(v)
}

func (s Transacted2) GetObjekteSha() schnittstellen.ShaLike {
	return s.ObjekteSha
}

func (s Transacted2) GetAkteSha() schnittstellen.ShaLike {
	return s.Metadatei.AkteSha
}

func (s *Transacted2) SetAkteSha(sh schnittstellen.ShaLike) {
	s.Metadatei.AkteSha = sha.Make(sh)
}

func (s Transacted2) GetTransactionIndex() values.Int {
	return s.TransactionIndex
}

func (o Transacted2) GetKey() string {
	return kennung.FormattedString(o.GetKennung())
}
