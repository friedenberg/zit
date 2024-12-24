package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_lock"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/object_id_provider"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// Saves the blob if necessary, applies the proto object, runs pre-commit hooks,
// runs the new hook, validates the blob, then calculates the sha for the object
func (s *Store) tryRealize(
	el sku.ExternalLike, mutter *sku.Transacted,
	o sku.CommitOptions,
) (err error) {
	if err = s.SaveBlob(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	kinder := el.GetSku()

	if mutter == nil && o.Contains(object_mode.ModeApplyProto) {
		s.protoZettel.Apply(kinder, kinder)
	}

	if genres.Type == el.GetSku().GetGenre() {
		if el.GetSku().GetType().IsEmpty() {
			el.GetSku().GetMetadata().Type = builtin_types.DefaultOrPanic(genres.Type)
		}
	}

	if err = s.tryPreCommitHooks(kinder, mutter, o); err != nil {
		if s.config.IgnoreHookErrors {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	// TODO just just mutter == nil
	if mutter == nil {
		if err = s.tryNewHook(kinder, o); err != nil {
			if s.config.IgnoreHookErrors {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}
	}

	if err = s.validate(kinder, mutter, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = kinder.CalculateObjectShas(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO add RealizeAndOrStore result
func (s *Store) tryRealizeAndOrStore(
	el sku.ExternalLike,
	o sku.CommitOptions,
) (err error) {
	child := el.GetSku()

	ui.Log().Printf("%s -> %s", o, child)

	if !s.GetDirectoryLayout().GetLockSmith().IsAcquired() &&
		o.ContainsAny(
			object_mode.ModeAddToInventoryList,
		) {
		err = errors.Wrap(file_lock.ErrLockRequired{
			Operation: "commit",
		})

		return
	}

	// TAI must be set before calculating object sha
	if o.ContainsAny(
		object_mode.ModeAddToInventoryList,
		object_mode.ModeUpdateTai,
	) {
		if o.Clock == nil {
			o.Clock = s
		}

		child.SetTai(o.Clock.GetTai())
	}

	if o.ContainsAny(
		object_mode.ModeAddToInventoryList,
	) && (child.ObjectId.IsEmpty() ||
		child.GetGenre() == genres.None ||
		child.GetGenre() == genres.Blob) {
		var ken *ids.ZettelId

		if ken, err = s.zettelIdIndex.CreateZettelId(); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = child.ObjectId.SetWithIdLike(ken); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	var parent *sku.Transacted

	if parent, err = s.fetchParentIfNecessary(child); err != nil {
		err = errors.Wrap(err)
		return
	}

	if parent != nil {
		defer sku.GetTransactedPool().Put(parent)
		child.Metadata.Cache.ParentTai = parent.GetTai()
	}

	if err = s.tryRealize(el, parent, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	if o.Contains(object_mode.ModeAddToInventoryList) {
		if err = s.addMissingTypeAndTags(o, child); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = s.addObjectToAbbrStore(child); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	// short circuits if the parent is equal to the child
	if o.Mode != object_mode.ModeReindex &&
		parent != nil &&
		ids.Equals(child.GetObjectId(), parent.GetObjectId()) &&
		child.Metadata.EqualsSansTai(&parent.Metadata) {

		sku.TransactedResetter.ResetWithExceptFields(child, parent)

		if o.Mode.Contains(object_mode.ModeLatest) {
			if err = s.ui.TransactedUnchanged(child); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		return
	}

	if err = s.GetConfig().ApplyDormantAndRealizeTags(
		child,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if o.Mode.Contains(object_mode.ModeLatest) {
		if err = s.GetConfig().AddTransacted(
			child,
			parent,
			s.GetBlobStore(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if child.GetGenre() == genres.Zettel {
			if err = s.zettelIdIndex.AddZettelId(&child.ObjectId); err != nil {
				if errors.Is(err, object_id_provider.ErrDoesNotExist{}) {
					ui.Log().Printf("object id does not contain value: %s", err)
					err = nil
				} else {
					err = errors.Wrapf(err, "failed to write zettel to index: %s", child)
					return
				}
			}
		}

	}

	if o.Contains(object_mode.ModeAddToInventoryList) {
		ui.Log().Print("adding to bestandsaufnahme", o, child)
		if err = s.commitTransacted(child, parent); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if o.Contains(object_mode.ModeLatest) {
		if err = s.GetStreamIndex().Add(
			child,
			child.GetObjectId().String(),
			o,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if parent == nil {
			if err = s.ui.TransactedNew(child); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else {
			// [are/kabuto !task project-2021-zit-features zz-inbox] add delta printing to changed objects
			// if err = s.Updated(mutter); err != nil {
			// 	err = errors.Wrap(err)
			// 	return
			// }

			if err = s.ui.TransactedUpdated(child); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	}

	if o.Contains(object_mode.ModeMergeCheckedOut) {
		if err = s.ReadExternalAndMergeIfNecessary(
			child,
			parent,
			o,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) fetchParentIfNecessary(
	sk *sku.Transacted,
) (mutter *sku.Transacted, err error) {
	mutter = sku.GetTransactedPool().Get()
	if err = s.GetStreamIndex().ReadOneObjectId(
		sk.GetObjectId().String(),
		mutter,
	); err != nil {
		if collections.IsErrNotFound(err) || errors.IsNotExist(err) {
			// TODO decide if this should continue to virtual stores
			sku.GetTransactedPool().Put(mutter)
			mutter = nil
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	sk.Metadata.Mutter().ResetWith(mutter.Metadata.Sha())

	return
}

// TODO add results for which stores had which change types
func (s *Store) commitTransacted(
	kinder *sku.Transacted,
	mutter *sku.Transacted,
) (err error) {
	sk := sku.GetTransactedPool().Get()

	sku.TransactedResetter.ResetWith(sk, kinder)

	if err = s.inventoryList.Add(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) handleUnchanged(
	t *sku.Transacted,
) (err error) {
	return s.ui.TransactedUnchanged(t)
}

func (s *Store) UpdateKonfig(
	sh interfaces.Sha,
) (kt *sku.Transacted, err error) {
	return s.CreateOrUpdateBlobSha(
		&ids.Config{},
		sh,
	)
}

func (s *Store) createTagsOrType(k *ids.ObjectId) (err error) {
	t := sku.GetTransactedPool().Get()
	defer sku.GetTransactedPool().Put(t)

	switch k.GetGenre() {
	default:
		err = genres.MakeErrUnsupportedGenre(k.GetGenre())
		return

	case genres.Type:
		t.GetMetadata().Type = builtin_types.DefaultOrPanic(genres.Type)

	case genres.Tag:
	}

	if err = t.ObjectId.SetWithIdLike(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.addObjectToAbbrStore(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.tryRealizeAndOrStore(
		t,
		sku.CommitOptions{Mode: object_mode.ModeCommit},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) addType(
	t ids.Type,
) (err error) {
	if t.IsEmpty() {
		err = errors.Errorf("attempting to add empty type")
		return
	}

	if err = s.GetStreamIndex().ObjectExists(t); err == nil {
		return
	}

	err = nil

	var k ids.ObjectId

	if err = k.SetWithIdLike(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.createTagsOrType(&k); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) addTypeAndExpandedIfNecessary(
	t1 ids.Type,
) (err error) {
	if t1.IsEmpty() {
		return
	}

	if builtin_types.IsBuiltin(t1) {
		return
	}

	typenExpanded := ids.ExpandOneSlice(
		t1,
		ids.MakeType,
		expansion.ExpanderRight,
	)

	for _, t := range typenExpanded {
		if err = s.addType(t); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) addTagAndExpanded(
	e ids.Tag,
) (err error) {
	if e.IsVirtual() {
		return
	}

	etikettenExpanded := ids.ExpandOneSlice(
		e,
		ids.MakeTag,
		expansion.ExpanderRight,
	)

	s.tagLock.Lock()
	defer s.tagLock.Unlock()

	for _, e1 := range etikettenExpanded {
		if e1.IsVirtual() {
			continue
		}

		if err = s.GetStreamIndex().ObjectExists(e1); err == nil {
			continue
		}

		err = nil

		var k ids.ObjectId

		if err = k.SetWithIdLike(e1); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = s.createTagsOrType(&k); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) addMissingTypeAndTags(
	co sku.CommitOptions,
	m *sku.Transacted,
) (err error) {
	t := m.GetType()

	if !co.DontAddMissingType {
		if err = s.addTypeAndExpandedIfNecessary(t); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if !co.DontAddMissingTags {
		es := quiter.SortedValues(m.Metadata.GetTags())

		for _, e := range es {
			if err = s.addTagAndExpanded(e); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

func (s *Store) addObjectToAbbrStore(m *sku.Transacted) (err error) {
	if err = s.GetAbbrStore().AddObjectToAbbreviationStore(m); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) reindexOne(besty, sk *sku.Transacted) (err error) {
	o := sku.CommitOptions{
		Mode: object_mode.ModeReindex,
	}

	if err = s.tryRealizeAndOrStore(sk, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetAbbrStore().AddObjectToAbbreviationStore(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
