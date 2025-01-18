package store_config

import (
	"encoding/gob"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/mutable_config_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
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

type (
	immutable_config_private = immutable_config.Config
	mutable_config_private   = mutable_config_blobs.Blob
	cli                      = config_mutable_cli.Config
	ApproximatedType         = blob_store.ApproximatedType

	Store interface {
		immutable_config.Config
		interfaces.Config

		ids.InlineTypeChecker
		GetTypeExtension(string) string
		GetCLIConfig() config_mutable_cli.Config
		GetImmutableConfig() immutable_config.Config
		GetMutableConfig() mutable_config_blobs.Blob
		GetFileExtensions() interfaces.FileExtensionGetter
		HasChanges() (ok bool)
		GetChanges() (out []string)

		GetTagOrRepoIdOrType(
			v string,
		) (sk *sku.Transacted, err error)
		GetImplicitTags(*ids.Tag) ids.TagSet
		GetApproximatedType(
			k interfaces.ObjectId,
		) (ct ApproximatedType)
		GetSku() *sku.Transacted
	}

	StoreMutable interface {
		Store

		AddTransacted(
			child *sku.Transacted,
			parent *sku.Transacted,
			ak *blob_store.VersionedStores,
		) (err error)

		Initialize(
			dirLayout env_repo.Env,
			kcli config_mutable_cli.Config,
			blobStore *blob_store.VersionedStores,
		) (err error)

		Reset() error

		Flush(
			dirLayout env_repo.Env,
			blobStore *blob_store.VersionedStores,
			printerHeader interfaces.FuncIter[string],
		) (err error)
	}
)

func Make() StoreMutable {
	return &env{}
}

type env struct {
	cli
	compiled
	immutable_config_private
}

func (a *env) GetCLIConfig() config_mutable_cli.Config {
	return a.cli
}

func (a *compiled) Reset() error {
	a.mutable_config_private = mutable_config_blobs.V1{}
	a.ExtensionsToTypes = make(map[string]string)
	a.TypesToExtensions = make(map[string]string)

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

func (a *env) GetMutableConfig() mutable_config_blobs.Blob {
	return a.mutable_config_private
}

func (c *env) Initialize(
	dirLayout env_repo.Env,
	kcli config_mutable_cli.Config,
	blobStore *blob_store.VersionedStores,
) (err error) {
	c.cli = kcli
	c.Reset()
	c.immutable_config_private = dirLayout.GetConfig()

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

func (kc *env) SetCli(k config_mutable_cli.Config) {
	kc.cli = k
}

func (kc *env) SetCliFromCommander(k config_mutable_cli.Config) {
	oldBasePath := kc.BasePath
	kc.cli = k
	kc.BasePath = oldBasePath
}

func (k *env) IsDryRun() bool {
	return k.DryRun
}

func (k *env) SetDryRun(v bool) {
	k.DryRun = v
}

func (k *env) GetTypeStringFromExtension(t string) string {
	return k.ExtensionsToTypes[t]
}

func (k *env) GetTypeExtension(v string) string {
	return k.TypesToExtensions[v]
}

func (k *env) AddTransacted(
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

func (kc *env) IsInlineType(k ids.Type) (isInline bool) {
	todo.Change("fix this horrible hack")
	if k.IsEmpty() {
		return true
	}

	isInline = kc.InlineTypes.ContainsKey(k.String()) ||
		builtin_types.IsBuiltin(k)

	return
}
