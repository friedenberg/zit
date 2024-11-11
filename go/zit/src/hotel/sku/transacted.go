package sku

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
)

type Transacted struct {
	ObjectId ids.ObjectId
	Metadata object_metadata.Metadata

	ExternalType ids.Type

	// TODO add support for querying the below
	RepoId           ids.RepoId
	State            external_state.State
	ExternalObjectId ids.ObjectId
}

func (t *Transacted) GetSkuExternal() *Transacted {
	return t
}

func (t *Transacted) GetRepoId() ids.RepoId {
	return t.RepoId
}

func (t *Transacted) GetExternalObjectId() ids.ExternalObjectId {
	return &t.ExternalObjectId
}

func (t *Transacted) GetExternalState() external_state.State {
	return t.State
}

func (a *Transacted) CloneTransacted() (b *Transacted) {
	b = GetTransactedPool().Get()
	TransactedResetter.ResetWith(b, a)
	return
}

func (a *Transacted) CloneExternalLike() ExternalLike {
	b := GetTransactedPool().Get()
	TransactedResetter.ResetWith(b, a)
	return b
}

func (t *Transacted) GetSku() *Transacted {
	return t
}

func (a *Transacted) SetFromTransacted(b *Transacted) (err error) {
	TransactedResetter.ResetWith(a, b)

	return
}

func (a *Transacted) Less(b *Transacted) bool {
	less := a.GetTai().Less(b.GetTai())

	return less
}

func (a *Transacted) GetTags() ids.TagSet {
	return a.Metadata.GetTags()
}

func (a *Transacted) AddTagPtr(e *ids.Tag) (err error) {
	if a.ObjectId.GetGenre() == genres.Tag &&
		strings.HasPrefix(a.ObjectId.String(), e.String()) {
		return
	}

	ek := a.Metadata.Cache.GetImplicitTags().KeyPtr(e)

	if a.Metadata.Cache.GetImplicitTags().ContainsKey(ek) {
		return
	}

	if err = a.GetMetadata().AddTagPtr(e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Transacted) AddTagPtrFast(e *ids.Tag) (err error) {
	if err = a.GetMetadata().AddTagPtrFast(e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Transacted) GetType() ids.Type {
	return a.Metadata.Type
}

func (a *Transacted) GetMetadata() *object_metadata.Metadata {
	return &a.Metadata
}

func (a *Transacted) GetTai() ids.Tai {
	return a.Metadata.GetTai()
}

func (a *Transacted) SetTai(t ids.Tai) {
	a.GetMetadata().Tai = t
}

func (a *Transacted) GetObjectId() *ids.ObjectId {
	return &a.ObjectId
}

func (a *Transacted) SetObjectIdLike(kl ids.IdLike) (err error) {
	if err = a.ObjectId.SetWithIdLike(kl); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Transacted) EqualsAny(b any) (ok bool) {
	return values.Equals(a, b)
}

func (a *Transacted) Equals(b *Transacted) (ok bool) {
	if a.GetObjectId().String() != b.GetObjectId().String() {
		return
	}

	// TODO-P2 determine why object shas in import test differed
	// if !a.Metadata.Sha().Equals(b.Metadata.Sha()) {
	// 	return
	// }

	if !a.Metadata.Equals(&b.Metadata) {
		return
	}

	return true
}

func (s *Transacted) GetGenre() interfaces.Genre {
	return s.ObjectId.GetGenre()
}

func (s *Transacted) IsNew() bool {
	return s.Metadata.Mutter().IsNull()
}

func (s *Transacted) CalculateObjectShaDebug() (err error) {
	return s.calculateObjectSha(true)
}

func (s *Transacted) CalculateObjectShas() (err error) {
	return s.calculateObjectSha(false)
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

func (s *Transacted) calculateObjectSha(debug bool) (err error) {
	f := object_inventory_format.GetShaForContext

	if debug {
		f = object_inventory_format.GetShaForContextDebug
	}

	wg := quiter.MakeErrorWaitGroupParallel()

	wg.Do(
		s.makeShaCalcFunc(
			f,
			object_inventory_format.Formats.MetadataObjectIdParent(),
			s.Metadata.Sha(),
		),
	)

	wg.Do(
		s.makeShaCalcFunc(
			f,
			object_inventory_format.Formats.Metadata(),
			&s.Metadata.SelfMetadata,
		),
	)

	wg.Do(
		s.makeShaCalcFunc(
			f,
			object_inventory_format.Formats.MetadataSansTai(),
			&s.Metadata.SelfMetadataWithoutTai,
		),
	)

	return wg.GetError()
}

func (s *Transacted) SetDormant(v bool) {
	s.Metadata.Cache.Dormant.SetBool(v)
}

func (s *Transacted) SetObjectSha(v interfaces.Sha) (err error) {
	return s.GetMetadata().Sha().SetShaLike(v)
}

func (s *Transacted) GetObjectSha() interfaces.Sha {
	return s.GetMetadata().Sha()
}

func (s *Transacted) GetBlobSha() interfaces.Sha {
	return &s.Metadata.Blob
}

func (s *Transacted) SetBlobSha(sh interfaces.Sha) error {
	return s.Metadata.Blob.SetShaLike(sh)
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
