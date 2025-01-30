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
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
)

type Getter interface {
	GetRepoLayout() Env
}

type Env struct {
	env_local.Env

	config_immutable_io.ConfigLoaded

	basePath              string
	readOnlyBlobStorePath string
	lockSmith             interfaces.LockSmith

	interfaces.DirectoryPaths

	local, remote blobStore

	CopyingBlobStore
	ObjectStore
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

	s.lockSmith = file_lock.New(s.FileLock())

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
		decoder := config_immutable_io.Coder{}

		if err = decoder.DecodeFromFile(
			&s.ConfigLoaded,
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
	if s.local, err = MakeBlobStoreFromLayout(*s); err != nil {
		errors.Wrap(err)
		return
	}

	s.CopyingBlobStore = MakeCopyingBlobStore(s.Env, s.local, s.remote)

	s.ObjectStore = ObjectStore{
		basePath: s.basePath,
		Config: env_dir.MakeConfigFromImmutableBlobConfig(
			s.ImmutableConfig.GetBlobStoreConfigImmutable(),
		),
		DirectoryPaths: s.DirectoryPaths,
		TemporaryFS:    s.GetTempLocal(),
	}

	return
}

func (a Env) GetEnv() env_ui.Env {
	return a.Env
}

func (a Env) SansObjectAge() (b Env) {
	b = a

	b.ObjectStore.Config = env_dir.MakeConfig(
		a.ObjectStore.Config.GetBlobCompression(),
		nil,
		a.ObjectStore.Config.GetLockInternalFiles(),
	)

	return
}

func (a Env) SansObjectCompression() (b Env) {
	b = a
	compression := config_immutable.CompressionTypeNone

	b.Config = env_dir.MakeConfig(
		&compression,
		a.GetBlobEncryption(),
		a.GetLockInternalFiles(),
	)

	return
}

func (s Env) GetConfig() config_immutable_io.ConfigLoaded {
	return s.ConfigLoaded
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
	if h.ConfigLoaded.ImmutableConfig == nil {
		return config_immutable.CurrentStoreVersion
	} else {
		return h.ConfigLoaded.ImmutableConfig.GetStoreVersion()
	}
}
