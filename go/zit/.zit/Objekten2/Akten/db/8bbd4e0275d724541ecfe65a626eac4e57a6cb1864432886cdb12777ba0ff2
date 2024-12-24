package config

import (
	"encoding/gob"
	"fmt"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/golf/mutable_config_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/dormant_index"
	"code.linenisgreat.com/zit/go/zit/src/juliett/blob_store"
)

var typeExpander expansion.Expander

func init() {
	typeExpander = expansion.MakeExpanderRight(`-`)

	gob.Register(
		collections_value.MakeMutableValueSet[values.String](
			nil,
		),
	)

	gob.Register(
		collections_value.MakeValueSet[values.String](
			nil,
		),
	)

	gob.Register(quiter.StringerKeyer[values.String]{})
	gob.Register(quiter.StringerKeyerPtr[ids.Type, *ids.Type]{})
}

type (
	immutable_config_private = immutable_config.Config
	mutable_config_private   = mutable_config_blobs.Blob
)

type Compiled struct {
	cli
	compiled
	immutable_config_private
	dormant *dormant_index.Index
}

func (a *compiled) Reset() error {
	a.mutable_config_private = mutable_config_blobs.V1{}
	a.ExtensionsToTypes = make(map[string]string)
	a.TypesToExtensions = make(map[string]string)

	a.lock = &sync.Mutex{}
	a.Tags = collections_value.MakeMutableValueSet[*tag](nil)
	a.InlineTypes = collections_value.MakeMutableValueSet[values.String](
		nil,
	)
	a.ImplicitTags = make(implicitTagMap)
	a.Repos = sku.MakeTransactedMutableSet()
	a.Types = sku.MakeTransactedMutableSet()

	sku.TransactedResetter.Reset(&a.Sku)

	return nil
}

func (a *Compiled) GetMutableConfig() mutable_config_blobs.Blob {
	return a.mutable_config_private
}

type cli = mutable_config_blobs.Cli

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

func (c *Compiled) Initialize(
	dirLayout dir_layout.DirLayout,
	kcli mutable_config_blobs.Cli,
	dormant *dormant_index.Index,
	blobStore *blob_store.VersionedStores,
) (err error) {
	c.cli = kcli
	c.Reset()
	c.immutable_config_private = dirLayout.GetConfig()
	c.dormant = dormant

	wg := quiter.MakeErrorWaitGroupParallel()
	wg.Do(func() (err error) {
		if err = c.loadMutableConfig(dirLayout, blobStore); err != nil {
			if errors.IsNotExist(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}
		}

		return
	})

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.cli.ApplyPrintOptionsConfig(c.GetPrintOptions())

	return
}

func (kc *Compiled) SetCli(k mutable_config_blobs.Cli) {
	kc.cli = k
}

func (kc *Compiled) SetCliFromCommander(k mutable_config_blobs.Cli) {
	oldBasePath := kc.BasePath
	kc.cli = k
	kc.BasePath = oldBasePath
}

func (kc *compiled) IsInlineType(k ids.Type) (isInline bool) {
	todo.Change("fix this horrible hack")
	if k.IsEmpty() {
		return true
	}

	isInline = kc.InlineTypes.ContainsKey(k.String()) ||
		builtin_types.IsBuiltin(k)

	return
}

type ApproximatedType = blob_store.ApproximatedType

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

	kag := blobStore.GetConfig()

	if k.mutable_config_private, _, err = kag.ParseTypedBlob(
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

func (k *Compiled) IsDryRun() bool {
	return k.DryRun
}

func (k *Compiled) SetDryRun(v bool) {
	k.DryRun = v
}

func (k *Compiled) GetTypeStringFromExtension(t string) string {
	return k.ExtensionsToTypes[t]
}

func (k *Compiled) GetTypeExtension(v string) string {
	return k.TypesToExtensions[v]
}

func (k *Compiled) AddTransacted(
	kinder *sku.Transacted,
	mutter *sku.Transacted,
	ak *blob_store.VersionedStores,
) (err error) {
	didChange := false

	g := kinder.ObjectId.GetGenre()

	switch g {
	case genres.Type:
		if didChange, err = k.addType(kinder); err != nil {
			err = errors.Wrap(err)
			return
		}

		if didChange {
			k.SetNeedsRecompile(fmt.Sprintf("modified type: %s", kinder))
		}

		return

	case genres.Tag:
		if didChange, err = k.addTag(kinder, mutter); err != nil {
			err = errors.Wrap(err)
			return
		}

	case genres.Repo:
		if didChange, err = k.addRepo(kinder); err != nil {
			err = errors.Wrap(err)
			return
		}

	case genres.Config:
		if didChange, err = k.setTransacted(kinder, ak); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if g != genres.Tag {
		return
	}

	if !didChange {
		return
	}

	if mutter == nil {
		return
	}

	if quiter.SetEquals(kinder.Metadata.Tags, mutter.Metadata.Tags) {
		return
	}

	k.SetNeedsRecompile(fmt.Sprintf("modified: %s", kinder))

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
