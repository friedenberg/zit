package env_workspace

import (
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_mutable_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/workspace_config_blobs"
)

const FileWorkspace = ".zit-workspace"

type Env interface {
	env_dir.Env
	InWorkspace() bool
	GetWorkspaceConfig() workspace_config_blobs.Blob
	GetDefaults() config_mutable_blobs.Defaults
	CreateWorkspace() (err error)
}

func Make(
	envDir env_dir.Env,
	configMutableBlob config_mutable_blobs.Blob,
) (out *env, err error) {
	out = &env{
		Env:           envDir,
		configMutable: configMutableBlob,
	}

	object := ids.TypeWithObject[*workspace_config_blobs.Blob]{
		Type: &ids.Type{},
	}

	if err = workspace_config_blobs.DecodeFromFile(
		&object,
		filepath.Join(out.GetCwd(), FileWorkspace),
	); errors.IsNotExist(err) {
		err = nil
	} else if err != nil {
		err = errors.Wrap(err)
		return
	} else {
		out.blob = *object.Object
	}

	defaults := out.configMutable.GetDefaults()

	out.defaults = config_mutable_blobs.DefaultsV1{
		Type: defaults.GetType(),
		Tags: defaults.GetTags(),
	}

	if out.blob != nil {
		defaults = out.blob.GetDefaults()

		if newType := defaults.GetType(); !newType.IsEmpty() {
			ui.Debug().Print(newType)
			out.defaults.Type = newType
		}

		if newTags := defaults.GetTags(); newTags.Len() > 0 {
			ui.Debug().Print(newTags)
			out.defaults.Tags = newTags
		}
	}

	return
}

type env struct {
	env_dir.Env
	configMutable config_mutable_blobs.Blob
	blob          workspace_config_blobs.Blob
	defaults      config_mutable_blobs.DefaultsV1
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

func (env *env) CreateWorkspace() (err error) {
	env.blob = &workspace_config_blobs.V0{}
	tipe := builtin_types.GetOrPanic(builtin_types.WorkspaceConfigTypeTomlV0).Type

	object := ids.TypeWithObject[*workspace_config_blobs.Blob]{
		Type:   &tipe,
		Object: &env.blob,
	}

	if err = workspace_config_blobs.EncodeToFile(
		&object,
		filepath.Join(env.GetCwd(), FileWorkspace),
	); errors.IsNotExist(err) {
		err = errors.Wrap(err)
		return
	}

	return
}
