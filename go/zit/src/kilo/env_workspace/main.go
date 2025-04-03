package env_workspace

import (
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/triple_hyphen_io"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_mutable_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/hotel/workspace_config_blobs"
)

const FileWorkspace = ".zit-workspace"

type Env interface {
	env_dir.Env
	GetWorkspaceDir() string
	AssertNotTemporary(errors.Context)
	AssertNotTemporaryOrOfferToCreate(errors.Context)
	IsTemporary() bool
	InWorkspace() bool
	GetWorkspaceConfig() workspace_config_blobs.Blob
	GetDefaults() config_mutable_blobs.Defaults
	CreateWorkspace(workspace_config_blobs.Blob) (err error)
	DeleteWorkspace() (err error)
}

func Make(
	envLocal env_local.Env,
	configMutableBlob config_mutable_blobs.Blob,
) (out *env, err error) {
	out = &env{
		Env:           envLocal,
		configMutable: configMutableBlob,
	}

	object := triple_hyphen_io.TypedStruct[*workspace_config_blobs.Blob]{
		Type: &ids.Type{},
	}

	if err = workspace_config_blobs.DecodeFromFile(
		&object,
		out.GetWorkspaceConfigFilePath(),
	); errors.IsNotExist(err) {
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

	return
}

type env struct {
	env_local.Env
	configMutable config_mutable_blobs.Blob
	blob          workspace_config_blobs.Blob
	defaults      config_mutable_blobs.DefaultsV1
}

func (env *env) GetWorkspaceDir() string {
	if env.IsTemporary() {
		// TODO return temp dir
		// return env.GetCwd()
	} else {
		return env.GetCwd()
	}

	return env.GetCwd()
}

func (env *env) GetWorkspaceConfigFilePath() string {
	return filepath.Join(env.GetWorkspaceDir(), FileWorkspace)
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
	return !env.InWorkspace()
}

func (env *env) InWorkspace() bool {
	return env.blob != nil
}

func (env *env) GetWorkspaceConfig() workspace_config_blobs.Blob {
	return env.blob
}

func (env *env) GetDefaults() config_mutable_blobs.Defaults {
	return env.defaults
}

func (env *env) CreateWorkspace(blob workspace_config_blobs.Blob) (err error) {
	env.blob = blob
	tipe := builtin_types.GetOrPanic(builtin_types.WorkspaceConfigTypeTomlV0).Type

	object := triple_hyphen_io.TypedStruct[*workspace_config_blobs.Blob]{
		Type:   &tipe,
		Struct: &env.blob,
	}

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
