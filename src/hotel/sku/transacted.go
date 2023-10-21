package sku

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type Transacted struct {
	Kennung          kennung.Kennung2
	Metadatei        metadatei.Metadatei
	ObjekteSha       sha.Sha
	TransactionIndex values.Int
	Kopf             kennung.Tai
}

func (t *Transacted) SetFromSkuLike(sk SkuLikePtr) (err error) {
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

func MakeSkuLikeSansObjekteSha(
	m metadatei.Metadatei,
	k kennung.Kennung,
) (sk *Transacted, err error) {
	sk = GetTransactedPool().Get()
	sk.Metadatei = m

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
	if sk, err = MakeSkuLikeSansObjekteSha(m, k); err != nil {
		err = errors.Wrap(err)
		return
	}

	sk.SetObjekteSha(os)

	return
}

func (a Transacted) String() string {
	return fmt.Sprintf(
		"%s %s %s",
		a.Kennung,
		a.ObjekteSha,
		a.Metadatei.AkteSha,
	)
}

func (a *Transacted) GetSkuLikePtr() SkuLikePtr {
	return a
}

func (a Transacted) GetEtiketten() kennung.EtikettSet {
	return a.Metadatei.GetEtiketten()
}

func (a Transacted) GetTyp() kennung.Typ {
	return a.Metadatei.Typ
}

func (a Transacted) GetMetadatei() metadatei.Metadatei {
	return a.Metadatei
}

func (a *Transacted) GetMetadateiPtr() *metadatei.Metadatei {
	return &a.Metadatei
}

func (a *Transacted) SetMetadatei(m metadatei.Metadatei) {
	a.Metadatei = m
}

func (a Transacted) GetTai() kennung.Tai {
	return a.GetMetadatei().GetTai()
}

func (a Transacted) GetKopf() kennung.Tai {
	return a.Kopf
}

func (a *Transacted) SetTai(t kennung.Tai) {
	a.GetMetadateiPtr().Tai = t
}

func (a Transacted) GetKennung() kennung.Kennung {
	return a.Kennung
}

func (a *Transacted) GetKennungPtr() kennung.KennungPtr {
	return &a.Kennung
}

func (a Transacted) GetKennungLike() kennung.Kennung {
	return a.Kennung
}

func (a *Transacted) GetKennungLikePtr() kennung.KennungPtr {
	return &a.Kennung
}

func (a *Transacted) SetKennungLike(kl kennung.Kennung) (err error) {
	if err = a.Kennung.SetWithKennung(kl); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Transacted) Reset() {
	a.Kopf.Reset()
	a.ObjekteSha.Reset()
	a.Kennung.SetGattung(gattung.Unknown)
	a.Metadatei.Reset()
	a.TransactionIndex.Reset()
}

func (a *Transacted) ResetWith(b Transacted) {
	a.Kopf = b.Kopf
	a.ObjekteSha = b.ObjekteSha
	a.Kennung.ResetWithKennung(b.Kennung)
	a.Metadatei.ResetWith(b.Metadatei)
	a.TransactionIndex.SetInt(b.TransactionIndex.Int())
}

// TODO-P2 switch this to default
func (a *Transacted) ResetWithPtr(
	b *Transacted,
) {
	a.ResetWith(*b)
}

func (a Transacted) Less(b Transacted) (ok bool) {
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

func (a Transacted) EqualsSkuLikePtr(b SkuLikePtr) bool {
	return values.Equals(a, b) || values.EqualsPtr(a, b)
}

func (a Transacted) EqualsAny(b any) (ok bool) {
	return values.Equals(a, b)
}

func (a Transacted) Equals(b Transacted) (ok bool) {
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

func (s Transacted) GetGattung() schnittstellen.GattungLike {
	return s.Kennung.GetGattung()
}

func (s *Transacted) IsNew() bool {
	return s.Metadatei.Verzeichnisse.Mutter.IsNull()
}

func (s *Transacted) SetObjekteSha(v schnittstellen.ShaLike) {
	s.ObjekteSha = sha.Make(v)
}

func (s Transacted) GetObjekteSha() schnittstellen.ShaLike {
	return s.ObjekteSha
}

func (s Transacted) GetAkteSha() schnittstellen.ShaLike {
	return s.Metadatei.AkteSha
}

func (s *Transacted) SetAkteSha(sh schnittstellen.ShaLike) {
	s.Metadatei.AkteSha = sha.Make(sh)
}

func (s Transacted) GetTransactionIndex() values.Int {
	return s.TransactionIndex
}

func (o Transacted) GetKey() string {
	return kennung.FormattedString(o.GetKennung())
}
