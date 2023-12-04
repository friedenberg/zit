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

type Transacted struct {
	Kennung          kennung.Kennung2
	Metadatei        metadatei.Metadatei
	ObjekteSha       sha.Sha
	TransactionIndex values.Int
	Kopf             kennung.Tai
}

func (t *Transacted) GetSkuLike() SkuLike {
	return t
}

func (t *Transacted) SetFromSkuLike(sk SkuLike) (err error) {
	if err = t.Kennung.SetWithKennung(sk.GetKennung()); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.ObjekteSha.SetShaLike(sk.GetObjekteSha())
	metadatei.Resetter.ResetWithPtr(&t.Metadatei, sk.GetMetadatei())
	t.GetMetadatei().Tai = sk.GetTai()

	t.Kopf = sk.GetTai()

	return
}

func (a *Transacted) Less(b *Transacted) bool {
	return a.GetTai().Less(b.GetTai())
}

func (a *Transacted) String() string {
	return fmt.Sprintf(
		"%s %s %s",
		&a.Kennung,
		&a.ObjekteSha,
		&a.Metadatei.AkteSha,
	)
}

func (a *Transacted) GetSkuLikePtr() SkuLike {
	return a
}

func (a *Transacted) GetEtiketten() kennung.EtikettSet {
	return a.Metadatei.GetEtiketten()
}

func (a *Transacted) GetTyp() kennung.Typ {
	return a.Metadatei.Typ
}

func (a *Transacted) GetMetadatei() *metadatei.Metadatei {
	return &a.Metadatei
}

func (a *Transacted) GetTai() kennung.Tai {
	return a.Metadatei.GetTai()
}

func (a *Transacted) GetKopf() kennung.Tai {
	return a.Kopf
}

func (a *Transacted) SetTai(t kennung.Tai) {
	a.GetMetadatei().Tai = t
}

func (a *Transacted) GetKennung() kennung.Kennung {
	return &a.Kennung
}

func (a *Transacted) GetKennungLike() kennung.Kennung {
	return &a.Kennung
}

func (a *Transacted) SetKennungLike(kl kennung.Kennung) (err error) {
	if err = a.Kennung.SetWithKennung(kl); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Transacted) EqualsSkuLikePtr(b SkuLike) bool {
	return values.Equals(a, b) || values.EqualsPtr(a, b)
}

func (a *Transacted) EqualsAny(b any) (ok bool) {
	return values.Equals(a, b)
}

func (a *Transacted) Equals(b *Transacted) (ok bool) {
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

	if !a.Metadatei.Equals(&b.Metadatei) {
		return
	}

	return true
}

func (s *Transacted) GetGattung() schnittstellen.GattungLike {
	return s.Kennung.GetGattung()
}

func (s *Transacted) IsNew() bool {
	return s.Metadatei.Verzeichnisse.Mutter.IsNull()
}

func (s *Transacted) SetObjekteSha(v schnittstellen.ShaLike) {
	s.ObjekteSha.SetShaLike(v)
}

func (s *Transacted) GetObjekteSha() schnittstellen.ShaLike {
	return &s.ObjekteSha
}

func (s *Transacted) GetAkteSha() schnittstellen.ShaLike {
	return &s.Metadatei.AkteSha
}

func (s *Transacted) SetAkteSha(sh schnittstellen.ShaLike) {
	s.Metadatei.AkteSha.SetShaLike(sh)
}

func (s *Transacted) GetTransactionIndex() values.Int {
	return s.TransactionIndex
}

func (o *Transacted) GetKey() string {
	return kennung.FormattedString(o.GetKennung())
}

type transactedLessor struct{}

func (transactedLessor) Less(a, b *Transacted) bool {
	return a.GetTai().Less(b.GetTai())
}

func (transactedLessor) LessPtr(a, b *Transacted) bool {
	return a.GetTai().Less(b.GetTai())
}

type transactedEqualer struct{}

func (transactedEqualer) Equals(a, b *Transacted) bool {
	return a.Equals(b)
}

func (transactedEqualer) EqualsPtr(a, b *Transacted) bool {
	return a.EqualsSkuLikePtr(b)
}
