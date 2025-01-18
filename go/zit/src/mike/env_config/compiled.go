package env_config

import (
	"fmt"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/mutable_config_blobs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/blob_store"
)

type (
	immutable_config_private = immutable_config.Config
	mutable_config_private   = mutable_config_blobs.Blob
	cli                      = config_mutable_cli.Config
	ApproximatedType         = blob_store.ApproximatedType
)

type compiled struct {
	lock sync.Locker

	changes []string

	Sku sku.Transacted

	mutable_config_private

	DefaultTags  ids.TagSet
	Tags         interfaces.MutableSetLike[*tag]
	ImplicitTags implicitTagMap

	// Typen
	ExtensionsToTypes map[string]string
	TypesToExtensions map[string]string
	DefaultType       sku.Transacted // deprecated
	Types             sku.TransactedMutableSet
	InlineTypes       interfaces.SetLike[values.String]

	// Kasten
	Repos sku.TransactedMutableSet
}

func (k *compiled) GetSku() *sku.Transacted {
	return &k.Sku
}

func (k *compiled) setTransacted(
	kt1 *sku.Transacted,
	blobStore *blob_store.VersionedStores,
) (didChange bool, err error) {
	if !sku.TransactedLessor.LessPtr(&k.Sku, kt1) {
		return
	}

	k.lock.Lock()
	defer k.lock.Unlock()

	didChange = true

	sku.Resetter.ResetWith(&k.Sku, kt1)

	k.setNeedsRecompile(fmt.Sprintf("updated konfig: %s", &k.Sku))

	if err = k.loadMutableConfigBlob(
		blobStore,
		k.Sku.GetType(),
		k.Sku.GetBlobSha(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k *compiled) loadMutableConfigBlob(
	blobStore *blob_store.VersionedStores,
	mutableConfigType ids.Type,
	blobSha interfaces.Sha,
) (err error) {
	// k.lock.Lock()
	// defer k.lock.Unlock()

	if k.mutable_config_private, _, err = blobStore.GetConfig().ParseTypedBlob(
		mutableConfigType,
		blobSha,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k *compiled) addRepo(
	c *sku.Transacted,
) (didChange bool, err error) {
	k.lock.Lock()
	defer k.lock.Unlock()

	b := sku.GetTransactedPool().Get()

	sku.Resetter.ResetWith(b, c)

	if didChange, err = quiter.AddOrReplaceIfGreater(
		k.Repos,
		b,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k *compiled) addType(
	b1 *sku.Transacted,
) (didChange bool, err error) {
	if err = genres.Type.AssertGenre(b1); err != nil {
		err = errors.Wrap(err)
		return
	}

	b := sku.GetTransactedPool().Get()

	sku.Resetter.ResetWith(b, b1)

	k.lock.Lock()
	defer k.lock.Unlock()

	if didChange, err = quiter.AddOrReplaceIfGreater(
		k.Types,
		b,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
