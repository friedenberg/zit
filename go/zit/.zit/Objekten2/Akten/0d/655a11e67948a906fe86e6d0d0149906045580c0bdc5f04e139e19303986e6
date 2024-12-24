package dir_layout

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_lock"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
)

const (
	EnvDir = "DIR_ZIT"
	EnvBin = "BIN_ZIT"
)

type Getter interface {
	GetDirectoryLayout() DirLayout
}

type DirLayout struct {
	Primitive
	basePath              string
	readOnlyBlobStorePath string
	lockSmith             interfaces.LockSmith
	age                   *age.Age
	immutable_config      immutable_config.Config

	interfaces.DirectoryPaths

	local, remote blobStore

	CopyingBlobStore
	ObjectStore

	TempLocal, TempOS TemporaryFS
}

func Make(
	o Options,
	primitive Primitive,
) (s DirLayout, err error) {
	s.Primitive = primitive
	s.age = &age.Age{}

	if o.BasePath == "" {
		o.BasePath = os.Getenv(EnvDir)
	}

	if o.BasePath == "" {
		if o.BasePath, err = os.Getwd(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	s.basePath = o.BasePath
	s.readOnlyBlobStorePath = o.GetReadOnlyBlobStorePath()

	var dp directoryPaths

	switch s.sv.GetInt() {
	case 6:
		s.xdg.Data = s.basePath
		dp = &directoryV0{}

	default:
		dp = &directoryV1{}
	}

	if err = dp.init(s.sv, s.xdg); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.DirectoryPaths = dp
	s.TempLocal.basePath = s.DirZit(fmt.Sprintf("tmp-%d", s.pid))

	// TODO add support for failing on pre-existing temp local
	// if files.Exists(s.TempLocal.basePath) {
	// 	err = MakeErrTempAlreadyExists(s.TempLocal.basePath)
	// 	return
	// }

	if err = s.MakeDir(s.TempLocal.basePath); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !o.PermitNoZitDirectory {
		if ok := files.Exists(s.DirZit()); !ok {
			err = errors.Wrap(ErrNotInZitDir{})
			return
		}
	}

	s.lockSmith = file_lock.New(s.FileLock())

	if s.execPath, err = os.Executable(); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO switch to useing MakeCommonEnv()
	{
		if err = os.Setenv(EnvDir, s.basePath); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = os.Setenv(EnvBin, s.execPath); err != nil {
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

	if err = s.loadKonfigAngeboren(); err != nil {
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

	s.CopyingBlobStore = MakeCopyingBlobStore(s.local, s.remote)

	s.ObjectStore = ObjectStore{
		basePath:         s.basePath,
		age:              s.age,
		immutable_config: s.immutable_config,
		DirectoryPaths:   s.DirectoryPaths,
		TemporaryFS:      s.TempLocal,
	}

	return
}

func (a DirLayout) SansObjectAge() (b DirLayout) {
	b = a
	b.age = nil
	b.ObjectStore.age = nil
	return
}

func (a DirLayout) SansObjectCompression() (b DirLayout) {
	b = a
	b.immutable_config.CompressionType = immutable_config.CompressionTypeNone
	b.ObjectStore.immutable_config.CompressionType = b.immutable_config.CompressionType
	return
}

func (s DirLayout) GetConfig() immutable_config.Config {
	return s.immutable_config
}

func (s *DirLayout) loadKonfigAngeboren() (err error) {
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

	if err = dec.Decode(&s.immutable_config); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s DirLayout) GetLockSmith() interfaces.LockSmith {
	return s.lockSmith
}

func (s *DirLayout) Age() *age.Age {
	return s.age
}

func (s DirLayout) Cwd() string {
	return s.cwd
}

func (s DirLayout) Executable() string {
	return s.execPath
}

func (s DirLayout) AbsFromCwdOrSame(p string) (p1 string) {
	var err error
	p1, err = filepath.Abs(p)
	if err != nil {
		p1 = p
	}

	return
}

func (s DirLayout) RelToCwdOrSame(p string) (p1 string) {
	var err error

	if p1, err = filepath.Rel(s.Cwd(), p); err != nil {
		p1 = p
	}

	return
}

func stringSliceJoin(s string, vs []string) []string {
	return append([]string{s}, vs...)
}

func (s DirLayout) MakeDir(d string) (err error) {
	if err = os.MkdirAll(d, os.ModeDir|0o755); err != nil {
		err = errors.Wrapf(err, "Dir: %q", d)
		return
	}

	return
}

func (s DirLayout) ResetCache() (err error) {
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

func (h DirLayout) DataFileStoreVersion() string {
	return filepath.Join(h.xdg.Data, "version")
}

func (h DirLayout) MakeCommonEnv() map[string]string {
	return map[string]string{
		"ZIT_BIN": h.Executable(),
		"ZIT_DIR": h.Dir(),
	}
}
