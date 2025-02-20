package env_repo

import (
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/config_immutable"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_lock"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_immutable_io"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
)

type Env struct {
	env_local.Env

	config config_immutable_io.ConfigLoadedPrivate

	basePath              string
	readOnlyBlobStorePath string
	lockSmith             interfaces.LockSmith

	interfaces.DirectoryPaths

	local, remote interfaces.LocalBlobStore

	blob_store.CopyingBlobStore
}

func Make(
	env env_local.Env,
	o Options,
) (s Env, err error) {
	s.Env = env
	if o.BasePath == "" {
		o.BasePath = os.Getenv(env_dir.EnvDir)
	}

	if o.BasePath == "" {
		if o.BasePath, err = os.Getwd(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	s.basePath = o.BasePath
	s.readOnlyBlobStorePath = o.GetReadOnlyBlobStorePath()

	dp := &directoryV1{}

	if err = dp.init(
		s.GetStoreVersion(),
		s.GetXDG(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.DirectoryPaths = dp

	// TODO add support for failing on pre-existing temp local
	// if files.Exists(s.TempLocal.basePath) {
	// 	err = MakeErrTempAlreadyExists(s.TempLocal.basePath)
	// 	return
	// }

	if !o.PermitNoZitDirectory {
		if ok := files.Exists(s.DirZit()); !ok {
			err = errors.Wrap(ErrNotInZitDir{})
			return
		}
	}

	if err = s.MakeDirPerms(0o700, s.GetXDG().GetXDGPaths()...); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.lockSmith = file_lock.New(env, s.FileLock())

	// TODO switch to useing MakeCommonEnv()
	{
		if err = os.Setenv(env_dir.EnvDir, s.basePath); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = os.Setenv(
			env_dir.EnvBin,
			s.GetExecPath(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	{
		decoder := config_immutable_io.CoderPrivate{}

		if err = decoder.DecodeFromFile(
			&s.config,
			s.FileConfigPermanent(),
		); err != nil {
			errors.Wrap(err)
			return
		}
	}

	if err = s.setupStores(); err != nil {
		errors.Wrap(err)
		return
	}

	return
}

func (s *Env) setupStores() (err error) {
	s.local = s.MakeBlobStore()
	s.CopyingBlobStore = blob_store.MakeCopyingBlobStore(s.Env, s.local, s.remote)

	return
}

func (a Env) GetEnv() env_ui.Env {
	return a.Env
}

func (s Env) GetConfigPublic() config_immutable_io.ConfigLoadedPublic {
	return config_immutable_io.ConfigLoadedPublic{
		Type:                     s.config.Type,
		ImmutableConfig:          s.config.ImmutableConfig.GetImmutableConfigPublic(),
		BlobStoreImmutableConfig: s.config.BlobStoreImmutableConfig,
	}
}

func (s Env) GetConfigPrivate() config_immutable_io.ConfigLoadedPrivate {
	return s.config
}

func (s Env) GetLockSmith() interfaces.LockSmith {
	return s.lockSmith
}

func stringSliceJoin(s string, vs []string) []string {
	return append([]string{s}, vs...)
}

func (s Env) ResetCache() (err error) {
	if err = files.SetAllowUserChangesRecursive(s.DirCache()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = os.RemoveAll(s.DirCache()); err != nil {
		err = errors.Wrapf(err, "failed to remove verzeichnisse dir")
		return
	}

	if err = s.MakeDir(s.DirCache()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.MakeDir(s.DirCacheObjects()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.MakeDir(s.DirCacheObjectPointers()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (h Env) DataFileStoreVersion() string {
	return filepath.Join(h.GetXDG().Data, "version")
}

func (h Env) GetStoreVersion() interfaces.StoreVersion {
	if h.config.ImmutableConfig == nil {
		return config_immutable.CurrentStoreVersion
	} else {
		return h.config.ImmutableConfig.GetStoreVersion()
	}
}

func (s Env) MakeBlobStore() interfaces.LocalBlobStore {
	return blob_store.MakeShardedFilesStore(
		s.DirBlobs(),
		env_dir.MakeConfigFromImmutableBlobConfig(
			s.GetConfigPublic().ImmutableConfig.GetBlobStoreConfigImmutable(),
		),
		s.GetTempLocal(),
	)
}
