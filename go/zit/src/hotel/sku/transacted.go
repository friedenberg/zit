package sku

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
)

type Transacted struct {
	Kennung          ids.ObjectId
	Metadatei        object_metadata.Metadata
	TransactionIndex values.Int
	Kopf             ids.Tai
}

func (t *Transacted) GetSkuLike() SkuLike {
	return t
}

func (a *Transacted) SetFromTransacted(b *Transacted) (err error) {
	TransactedResetter.ResetWith(a, b)

	return
}

func (t *Transacted) SetFromSkuLike(sk SkuLike) (err error) {
	if err = t.Kennung.SetWithIdLike(sk.GetObjectId()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = t.SetObjectSha(sk.GetObjekteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	object_metadata.Resetter.ResetWith(&t.Metadatei, sk.GetMetadata())
	t.GetMetadata().Tai = sk.GetTai()

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

func (a *Transacted) StringKennungBezeichnung() string {
	return fmt.Sprintf(
		"[%s %q]",
		&a.Kennung,
		a.Metadatei.Description,
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
		a.GetMetadata().Sha(),
	)
}

func (a *Transacted) StringKennungMutter() string {
	return fmt.Sprintf(
		"%s^@%s",
		&a.Kennung,
		a.GetMetadata().Mutter(),
	)
}

func (a *Transacted) GetSkuLikePtr() SkuLike {
	return a
}

func (a *Transacted) GetEtiketten() ids.TagSet {
	return a.Metadatei.GetTags()
}

func (a *Transacted) AddEtikettPtr(e *ids.Tag) (err error) {
	if a.Kennung.GetGenre() == genres.Tag {
		e1 := ids.MustTag(a.Kennung.String())
		ex := ids.ExpandOne(&e1, expansion.ExpanderRight)

		if ex.ContainsKey(ex.KeyPtr(e)) {
			return
		}
	}

	ek := a.Metadatei.Cache.GetImplicitTags().KeyPtr(e)

	if a.Metadatei.Cache.GetImplicitTags().ContainsKey(ek) {
		return
	}

	if err = a.GetMetadata().AddTagPtr(e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Transacted) AddEtikettPtrFast(e *ids.Tag) (err error) {
	if err = a.GetMetadata().AddTagPtrFast(e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Transacted) GetType() ids.Type {
	return a.Metadatei.Type
}

func (a *Transacted) GetMetadata() *object_metadata.Metadata {
	return &a.Metadatei
}

func (a *Transacted) GetTai() ids.Tai {
	return a.Metadatei.GetTai()
}

func (a *Transacted) GetKopf() ids.Tai {
	return a.Kopf
}

func (a *Transacted) SetTai(t ids.Tai) {
	// log.Debug().Caller(6, "before: %s", a.StringKennungTai())
	a.GetMetadata().Tai = t
	// log.Debug().Caller(6, "after: %s", a.StringKennungTai())
}

func (a *Transacted) GetObjectId() *ids.ObjectId {
	return &a.Kennung
}

func (a *Transacted) GetKennungLike() ids.IdLike {
	return &a.Kennung
}

func (a *Transacted) SetObjectIdLike(kl ids.IdLike) (err error) {
	if err = a.Kennung.SetWithIdLike(kl); err != nil {
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

	if a.GetObjectId().String() != b.GetObjectId().String() {
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

func (s *Transacted) GetGenre() interfaces.Genre {
	return s.Kennung.GetGenre()
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
	f func(object_inventory_format.FormatGeneric, object_inventory_format.FormatterContext) (*sha.Sha, error),
	of object_inventory_format.FormatGeneric,
	sh *sha.Sha,
) interfaces.FuncError {
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
	f := object_inventory_format.GetShaForContext

	if debug {
		f = object_inventory_format.GetShaForContextDebug
	}

	wg := iter.MakeErrorWaitGroupParallel()

	wg.Do(
		s.makeShaCalcFunc(
			f,
			object_inventory_format.Formats.MetadateiKennungMutter(),
			s.Metadatei.Sha(),
		),
	)

	wg.Do(
		s.makeShaCalcFunc(
			f,
			object_inventory_format.Formats.Metadatei(),
			&s.Metadatei.SelfMetadata,
		),
	)

	wg.Do(
		s.makeShaCalcFunc(
			f,
			object_inventory_format.Formats.MetadateiSansTai(),
			&s.Metadatei.SelfMetadataWithoutTai,
		),
	)

	return wg.GetError()
}

func (s *Transacted) SetSchlummernd(v bool) {
	s.Metadatei.Cache.Dormant.SetBool(v)
}

func (s *Transacted) SetObjectSha(v interfaces.Sha) (err error) {
	return s.GetMetadata().Sha().SetShaLike(v)
}

func (s *Transacted) GetObjekteSha() interfaces.Sha {
	return s.GetMetadata().Sha()
}

func (s *Transacted) GetAkteSha() interfaces.Sha {
	return &s.Metadatei.Blob
}

func (s *Transacted) SetBlobSha(sh interfaces.Sha) error {
	return s.Metadatei.Blob.SetShaLike(sh)
}

func (s *Transacted) GetTransactionIndex() values.Int {
	return s.TransactionIndex
}

func (o *Transacted) GetKey() string {
	return ids.FormattedString(o.GetObjectId())
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
