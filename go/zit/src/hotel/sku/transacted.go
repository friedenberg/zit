package sku

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/objekte_format"
)

type Transacted struct {
	Kennung          kennung.Kennung2
	Metadatei        metadatei.Metadatei
	TransactionIndex values.Int
	Kopf             kennung.Tai
}

func (t *Transacted) GetSkuLike() SkuLike {
	return t
}

func (a *Transacted) SetFromTransacted(b *Transacted) (err error) {
	TransactedResetter.ResetWith(a, b)

	return
}

func (t *Transacted) SetFromSkuLike(sk SkuLike) (err error) {
	if err = t.Kennung.SetWithKennung(sk.GetKennung()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = t.SetObjekteSha(sk.GetObjekteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	metadatei.Resetter.ResetWith(&t.Metadatei, sk.GetMetadatei())
	t.GetMetadatei().Tai = sk.GetTai()

	t.Kopf = sk.GetTai()

	return
}

func (a *Transacted) Less(b *Transacted) bool {
	less := a.GetTai().Less(b.GetTai())

	// 	op := ">"

	// 	if less {
	// 		op = "<"
	// 	}

	// 	log.Debug().Print(a.StringKennungTaiAkte(), op, b.StringKennungTaiAkte())

	return less
}

func (a *Transacted) String() string {
	return fmt.Sprintf(
		"%s %s %s",
		&a.Kennung,
		a.GetObjekteSha(),
		a.GetAkteSha(),
	)
}

func (a *Transacted) StringKennungTai() string {
	return fmt.Sprintf(
		"%s@%s",
		&a.Kennung,
		a.GetTai().StringDefaultFormat(),
	)
}

func (a *Transacted) StringKennungTaiAkte() string {
	return fmt.Sprintf(
		"%s@%s@%s",
		&a.Kennung,
		a.GetTai().StringDefaultFormat(),
		a.GetAkteSha(),
	)
}

func (a *Transacted) StringKennungSha() string {
	return fmt.Sprintf(
		"%s@%s",
		&a.Kennung,
		a.GetMetadatei().Sha(),
	)
}

func (a *Transacted) StringKennungMutter() string {
	return fmt.Sprintf(
		"%s^@%s",
		&a.Kennung,
		a.GetMetadatei().Mutter(),
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
	// log.Debug().Caller(6, "before: %s", a.StringKennungTai())
	a.GetMetadatei().Tai = t
	// log.Debug().Caller(6, "after: %s", a.StringKennungTai())
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
	return s.Metadatei.Mutter().IsNull()
}

func (s *Transacted) CalculateObjekteShaDebug() (err error) {
	return s.calculateObjekteSha(true)
}

func (s *Transacted) CalculateObjekteShas() (err error) {
	return s.calculateObjekteSha(false)
}

func (s *Transacted) makeShaCalcFunc(
	f func(objekte_format.FormatGeneric, objekte_format.FormatterContext) (*sha.Sha, error),
	of objekte_format.FormatGeneric,
	sh *sha.Sha,
) schnittstellen.FuncError {
	return func() (err error) {
		var actual *sha.Sha

		if actual, err = f(
			of,
			s,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer sha.GetPool().Put(actual)

		sh.ResetWith(actual)

		return
	}
}

func (s *Transacted) calculateObjekteSha(debug bool) (err error) {
	f := objekte_format.GetShaForContext

	if debug {
		f = objekte_format.GetShaForContextDebug
	}

	wg := iter.MakeErrorWaitGroupParallel()

	wg.Do(
		s.makeShaCalcFunc(
			f,
			objekte_format.Formats.MetadateiKennungMutter(),
			s.Metadatei.Sha(),
		),
	)

	wg.Do(
		s.makeShaCalcFunc(
			f,
			objekte_format.Formats.Metadatei(),
			&s.Metadatei.SelbstMetadatei,
		),
	)

	wg.Do(
		s.makeShaCalcFunc(
			f,
			objekte_format.Formats.MetadateiSansTai(),
			&s.Metadatei.SelbstMetadateiSansTai,
		),
	)

	return wg.GetError()
}

func (s *Transacted) SetArchiviert(v bool) {
	s.Metadatei.Verzeichnisse.Archiviert.SetBool(v)
}

func (s *Transacted) SetObjekteSha(v schnittstellen.ShaLike) (err error) {
	return s.GetMetadatei().Sha().SetShaLike(v)
}

func (s *Transacted) GetObjekteSha() schnittstellen.ShaLike {
	return s.GetMetadatei().Sha()
}

func (s *Transacted) GetAkteSha() schnittstellen.ShaLike {
	return &s.Metadatei.Akte
}

func (s *Transacted) SetAkteSha(sh schnittstellen.ShaLike) error {
	return s.Metadatei.Akte.SetShaLike(sh)
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
