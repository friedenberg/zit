package config

import (
	"encoding/gob"
	"fmt"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/mutable_config_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/dormant_index"
	"code.linenisgreat.com/zit/go/zit/src/lima/blob_store"
)

func init() {
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

func (c *Compiled) Initialize(
	dirLayout repo_layout.Layout,
	kcli config_mutable_cli.Config,
	dormant *dormant_index.Index,
	blobStore *blob_store.VersionedStores,
) (err error) {
	c.cli = kcli
	c.Reset()
	c.immutable_config_private = dirLayout.GetConfig()
	c.dormant = dormant

	wg := errors.MakeWaitGroupParallel()
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

func (kc *Compiled) SetCli(k config_mutable_cli.Config) {
	kc.cli = k
}

func (kc *Compiled) SetCliFromCommander(k config_mutable_cli.Config) {
	oldBasePath := kc.BasePath
	kc.cli = k
	kc.BasePath = oldBasePath
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
	child *sku.Transacted,
	parent *sku.Transacted,
	ak *blob_store.VersionedStores,
) (err error) {
	didChange := false

	g := child.ObjectId.GetGenre()

	switch g {
	case genres.Type:
		if didChange, err = k.addType(child); err != nil {
			err = errors.Wrap(err)
			return
		}

		if didChange {
			k.SetNeedsRecompile(fmt.Sprintf("modified type: %s", child))
		}

		return

	case genres.Tag:
		if didChange, err = k.addTag(child, parent); err != nil {
			err = errors.Wrap(err)
			return
		}

		var tag ids.Tag

		if err = tag.TodoSetFromObjectId(child.GetObjectId()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if child.Metadata.GetTags().Len() > 0 {
			k.SetNeedsRecompile(
				fmt.Sprintf(
					"tag with tags added: %q -> %q",
					tag,
					quiter.SortedValues(child.Metadata.GetTags()),
				),
			)
		}

	case genres.Repo:
		if didChange, err = k.addRepo(child); err != nil {
			err = errors.Wrap(err)
			return
		}

	case genres.Config:
		if didChange, err = k.setTransacted(child, ak); err != nil {
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

	if parent == nil {
		return
	}

	if quiter.SetEquals(child.Metadata.Tags, parent.Metadata.Tags) {
		return
	}

	k.SetNeedsRecompile(fmt.Sprintf("modified: %s", child))

	return
}

func (kc *Compiled) IsInlineType(k ids.Type) (isInline bool) {
	todo.Change("fix this horrible hack")
	if k.IsEmpty() {
		return true
	}

	isInline = kc.InlineTypes.ContainsKey(k.String()) ||
		builtin_types.IsBuiltin(k)

	return
}
