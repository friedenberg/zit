package repo_layout

import (
	"encoding/gob"
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
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

	config

	basePath              string
	readOnlyBlobStorePath string
	lockSmith             interfaces.LockSmith
	age                   *age.Age

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
	s.age = &age.Age{}

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

	if err = s.storeVersion.ReadFromFile(
		s.DataFileStoreVersion(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var dp directoryPaths

	switch s.GetStoreVersion().GetInt() {
	case 6:
		xdg := s.GetXDG()
		xdg.Data = s.basePath
		s.SetXDG(xdg)
		dp = &directoryV0{}

	default:
		dp = &directoryV1{}
	}

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

	{
		fa := s.FileAge()

		if files.Exists(fa) {
			var i age.Identity

			if err = i.SetFromPath(fa); err != nil {
				errors.Wrap(err)
				return
			}

			if err = s.age.AddIdentity(i); err != nil {
				errors.Wrap(err)
				return
			}
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

	if s.readOnlyBlobStorePath != "" {
		// ui.Err().Printf("using remote store: %q", s.readOnlyBlobStorePath)
		s.remote = MakeBlobStore(
			s.readOnlyBlobStorePath,
			s.age,
			immutable_config.CompressionTypeZstd,
		)
	}

	s.CopyingBlobStore = MakeCopyingBlobStore(s.Env, s.local, s.remote)

	s.ObjectStore = ObjectStore{
		basePath:       s.basePath,
		age:            s.age,
		config:         s.config,
		DirectoryPaths: s.DirectoryPaths,
		TemporaryFS:    s.GetDirLayout().TempLocal,
	}

	return
}

func (a Layout) SansObjectAge() (b Layout) {
	b = a
	b.age = nil
	b.ObjectStore.age = nil
	return
}

func (a Layout) SansObjectCompression() (b Layout) {
	b = a
	b.compressionType = immutable_config.CompressionTypeNone
	b.ObjectStore.config.compressionType = b.config.GetCompressionType()
	return
}

func (s Layout) GetConfig() immutable_config.Config {
	return s.config.Config
}

func (s *Layout) loadImmutableConfig() (err error) {
	var config immutable_config.Latest
	s.config.Config = &config

	var f *os.File

	if f, err = files.OpenExclusiveReadOnly(s.FileConfigPermanent()); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.DeferredCloser(&err, f)

	dec := gob.NewDecoder(f)

	// TODO use text object format

	if err = dec.Decode(&config); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.config.storeVersion = immutable_config.MakeStoreVersion(config.GetStoreVersion())
	s.config.compressionType = config.GetCompressionType()
	s.config.lockInternalFiles = config.GetLockInternalFiles()

	return
}

func (s Layout) GetLockSmith() interfaces.LockSmith {
	return s.lockSmith
}

func (s *Layout) Age() *age.Age {
	return s.age
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

func (h Layout) GetStoreVersion() immutable_config.StoreVersion {
	return h.storeVersion
}
