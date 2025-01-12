package repo_layout

import (
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_lock"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
)

type Getter interface {
	GetRepoLayout() Layout
}

type Layout struct {
	*env.Env

	Config

	basePath              string
	readOnlyBlobStorePath string
	lockSmith             interfaces.LockSmith

	interfaces.DirectoryPaths

	local, remote blobStore

	CopyingBlobStore
	ObjectStore
}

func Make(
	env *env.Env,
	o Options,
) (s Layout, err error) {
	s.Env = env

	if o.BasePath == "" {
		o.BasePath = os.Getenv(dir_layout.EnvDir)
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

	if err = dp.init(s.GetStoreVersion(), s.GetXDG()); err != nil {
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

	s.lockSmith = file_lock.New(s.FileLock())

	// TODO switch to useing MakeCommonEnv()
	{
		if err = os.Setenv(dir_layout.EnvDir, s.basePath); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = os.Setenv(
			dir_layout.EnvBin,
			s.GetExecPath(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.loadImmutableConfig(); err != nil {
		errors.Wrap(err)
		return
	}

	if s.local, err = MakeBlobStoreFromHome(s); err != nil {
		errors.Wrap(err)
		return
	}

	s.CopyingBlobStore = MakeCopyingBlobStore(s.Env, s.local, s.remote)

	s.ObjectStore = ObjectStore{
		basePath:       s.basePath,
		Config:         dir_layout.MakeConfigFromImmutableBlobConfig(s.Config.Config.GetBlobStoreImmutableConfig()),
		DirectoryPaths: s.DirectoryPaths,
		TemporaryFS:    s.GetDirLayout().TempLocal,
	}

	return
}

func (a Layout) SansObjectAge() (b Layout) {
	b = a

	b.ObjectStore.Config = dir_layout.MakeConfig(
		nil,
		a.ObjectStore.Config.GetCompressionType(),
		a.ObjectStore.Config.GetLockInternalFiles(),
	)

	return
}

func (a Layout) SansObjectCompression() (b Layout) {
	b = a

	b.ObjectStore.Config = dir_layout.MakeConfig(
		a.ObjectStore.Config.GetAgeEncryption(),
		immutable_config.CompressionTypeNone,
		a.ObjectStore.Config.GetLockInternalFiles(),
	)

	return
}

func (s Layout) GetConfig() immutable_config.Config {
	return s.Config.Config
}

func (s Layout) GetLockSmith() interfaces.LockSmith {
	return s.lockSmith
}

func stringSliceJoin(s string, vs []string) []string {
	return append([]string{s}, vs...)
}

func (s Layout) ResetCache() (err error) {
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

func (h Layout) DataFileStoreVersion() string {
	return filepath.Join(h.GetXDG().Data, "version")
}

func (h Layout) GetStoreVersion() interfaces.StoreVersion {
	if h.Config.Config == nil {
		return immutable_config.CurrentStoreVersion
	} else {
		return h.Config.Config.GetStoreVersion()
	}
}
