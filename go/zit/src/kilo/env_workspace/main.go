package env_workspace

import (
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/triple_hyphen_io"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_mutable_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/workspace_config_blobs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
)

type Env interface {
	env_dir.Env
	GetWorkspaceDir() string
	AssertNotTemporary(errors.Context)
	AssertNotTemporaryOrOfferToCreate(errors.Context)
	IsTemporary() bool
	GetWorkspaceConfig() workspace_config_blobs.Blob
	GetDefaults() config_mutable_blobs.Defaults
	CreateWorkspace(workspace_config_blobs.Blob) (err error)
	DeleteWorkspace() (err error)
	GetStore() *Store
	GetStoreFS() *store_fs.Store
}

func Make(
	envLocal env_local.Env,
	configMutableBlob config_mutable_blobs.Blob,
	skuConfig sku.Config,
	deletedPrinter interfaces.FuncIter[*fd.FD],
	fileExtensions interfaces.FileExtensionGetter,
	envRepo env_repo.Env,
) (out *env, err error) {
	out = &env{
		Env:           envLocal,
		configMutable: configMutableBlob,
	}

	object := triple_hyphen_io.TypedStruct[*workspace_config_blobs.Blob]{
		Type: &ids.Type{},
	}

	expectedWorkspaceConfigFilePath := filepath.Join(
		out.GetCwd(),
		env_repo.FileWorkspace,
	)

	if err = workspace_config_blobs.DecodeFromFile(
		&object,
		expectedWorkspaceConfigFilePath,
	); errors.IsNotExist(err) {
		out.isTemporary = true
		err = nil
	} else if err != nil {
		err = errors.Wrap(err)
		return
	} else {
		out.blob = *object.Struct
	}

	defaults := out.configMutable.GetDefaults()

	out.defaults = config_mutable_blobs.DefaultsV1{
		Type: defaults.GetType(),
		Tags: defaults.GetTags(),
	}

	if out.blob != nil {
		defaults = out.blob.GetDefaults()

		if newType := defaults.GetType(); !newType.IsEmpty() {
			out.defaults.Type = newType
		}

		if newTags := defaults.GetTags(); newTags.Len() > 0 {
			out.defaults.Tags = append(out.defaults.Tags, newTags...)
		}
	}

	if out.isTemporary {
		if out.dir, err = out.GetTempLocal().DirTempWithTemplate(
			"workspace-*",
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		out.dir = out.GetCwd()
	}

	if out.storeFS, err = store_fs.Make(
		skuConfig,
		deletedPrinter,
		fileExtensions,
		envRepo,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	out.store.StoreLike = out.storeFS

	return
}

type env struct {
	env_local.Env

	isTemporary bool

	// dir is populated on init to either the cwd, or a temporary directory,
	// depending on whether $PWD/.zit-workspace exists.
	//
	// Later, dir may be set to $PWD/.zit-workspace by CreateWorkspace
	dir string

	configMutable config_mutable_blobs.Blob
	blob          workspace_config_blobs.Blob
	defaults      config_mutable_blobs.DefaultsV1

	storeFS *store_fs.Store
	store   Store
}

func (env *env) GetWorkspaceDir() string {
	return env.dir
}

func (env *env) GetWorkspaceConfigFilePath() string {
	return filepath.Join(env.GetWorkspaceDir(), env_repo.FileWorkspace)
}

func (env *env) AssertNotTemporary(context errors.Context) {
	if env.IsTemporary() {
		context.CancelWithError(ErrNotInWorkspace{env: env})
	}
}

func (env *env) AssertNotTemporaryOrOfferToCreate(context errors.Context) {
	if env.IsTemporary() {
		context.CancelWithError(
			ErrNotInWorkspace{
				env:           env,
				offerToCreate: true,
			},
		)
	}
}

func (env *env) IsTemporary() bool {
	return env.isTemporary
}

func (env *env) GetWorkspaceConfig() workspace_config_blobs.Blob {
	return env.blob
}

func (env *env) GetDefaults() config_mutable_blobs.Defaults {
	return env.defaults
}

func (env *env) GetStore() *Store {
	return &env.store
}

func (env *env) GetStoreFS() *store_fs.Store {
	return env.storeFS
}

func (env *env) CreateWorkspace(blob workspace_config_blobs.Blob) (err error) {
	env.blob = blob
	tipe := builtin_types.GetOrPanic(builtin_types.WorkspaceConfigTypeTomlV0).Type

	object := triple_hyphen_io.TypedStruct[*workspace_config_blobs.Blob]{
		Type:   &tipe,
		Struct: &env.blob,
	}

	env.dir = env.GetCwd()

	if err = workspace_config_blobs.EncodeToFile(
		&object,
		env.GetWorkspaceConfigFilePath(),
	); errors.IsExist(err) {
		err = errors.BadRequestf("workspace already exists")
		return
	} else if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (env *env) DeleteWorkspace() (err error) {
	if err = env.Delete(env.GetWorkspaceConfigFilePath()); errors.IsNotExist(err) {
		err = nil
	} else if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
