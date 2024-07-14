package config

import (
	"encoding/gob"
	"fmt"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	pkg_angeboren "code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/mutable_config"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/india/dormant_index"
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

	gob.Register(iter.StringerKeyer[values.String]{})
	gob.Register(iter.StringerKeyerPtr[ids.Type, *ids.Type]{})
}

type immutable_config = pkg_angeboren.Config

type Compiled struct {
	cli
	compiled
	immutable_config
	dormant *dormant_index.Index
}

func (a *compiled) Reset() error {
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

func (a *Compiled) GetMutableConfig() *mutable_config.Blob {
	return &a.Blob
}

type cli = mutable_config.Cli

type compiled struct {
	lock sync.Locker

	changes []string

	Sku sku.Transacted

	mutable_config.Blob

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
	s fs_home.Home,
	kcli mutable_config.Cli,
	dormant *dormant_index.Index,
) (err error) {
	c.cli = kcli
	c.Reset()
	c.immutable_config = s.GetConfig()
	c.dormant = dormant

	wg := iter.MakeErrorWaitGroupParallel()
	wg.Do(func() (err error) {
		if err = c.loadMutableConfig(s); err != nil {
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

	return
}

func (kc *Compiled) SetCli(k mutable_config.Cli) {
	kc.cli = k
}

func (kc *Compiled) SetCliFromCommander(k mutable_config.Cli) {
	oldBasePath := kc.BasePath
	kc.cli = k
	kc.BasePath = oldBasePath
}

func (kc *compiled) IsInlineType(k ids.Type) (isInline bool) {
	todo.Change("fix this horrible hack")
	if k.IsEmpty() {
		return true
	}

	isInline = kc.InlineTypes.ContainsKey(k.String())

	return
}

type ApproximatedType = blob_store.ApproximatedType

func (k *compiled) setTransacted(
	kt1 *sku.Transacted,
	kag interfaces.BlobGetter[*mutable_config.Blob],
) (didChange bool, err error) {
	if !sku.TransactedLessor.LessPtr(&k.Sku, kt1) {
		return
	}

	k.lock.Lock()
	defer k.lock.Unlock()

	didChange = true

	if err = k.Sku.SetFromSkuLike(kt1); err != nil {
		err = errors.Wrap(err)
		return
	}

	k.setHasChanges(fmt.Sprintf("updated konfig: %s", &k.Sku))

	var a *mutable_config.Blob

	if a, err = kag.GetBlob(k.Sku.GetBlobSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	k.Blob = *a

	return
}

func (k *compiled) addRepo(
	c *sku.Transacted,
) (didChange bool, err error) {
	k.lock.Lock()
	defer k.lock.Unlock()

	b := sku.GetTransactedPool().Get()

	if err = b.SetFromSkuLike(c); err != nil {
		err = errors.Wrap(err)
		return
	}

	if didChange, err = iter.AddOrReplaceIfGreater(
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
	mode objekte_mode.Mode,
) (err error) {
	didChange := false

	switch kinder.ObjectId.GetGenre() {
	case genres.Type:
		if didChange, err = k.addType(kinder); err != nil {
			err = errors.Wrap(err)
			return
		}

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
		if didChange, err = k.setTransacted(kinder, ak.GetConfigV0()); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	// switch kinder.Kennung.GetGattung() {
	// case gattung.Typ, gattung.Etikett, gattung.Kasten:
	// 	didChange = didChange && mutter != nil
	// }

	if didChange && (mutter != nil || mode.Contains(objekte_mode.ModeSchwanz)) {
		k.SetHasChanges(fmt.Sprintf("added: %s", kinder))
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

	if err = b.SetFromSkuLike(b1); err != nil {
		err = errors.Wrap(err)
		return
	}

	k.lock.Lock()
	defer k.lock.Unlock()

	if didChange, err = iter.AddOrReplaceIfGreater(
		k.Types,
		b,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
