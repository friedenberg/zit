package env_workspace

import (
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/workspace_config_blobs"
)

type Env interface {
	GetWorkspaceConfig() (workspace_config_blobs.Blob, error)
}

func Make(envDir env_dir.Env) Env {
	return &env{Env: envDir}
}

type env struct {
	env_dir.Env
}

func (env *env) GetWorkspaceConfig() (blob workspace_config_blobs.Blob, err error) {
	object := ids.TypeWithObject[*workspace_config_blobs.Blob]{
		Type: &ids.Type{},
	}

	if err = workspace_config_blobs.DecodeFromFile(
		&object,
		filepath.Join(env.GetCwd(), ".zit-workspace.toml"),
	); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	blob = *object.Object

	return
}
