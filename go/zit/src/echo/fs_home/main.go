package fs_home

import (
	"encoding/gob"
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_lock"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
)

const (
	EnvDir = "DIR_ZIT"
	EnvBin = "BIN_ZIT"
)

type Getter interface {
	GetFSHome() Home
}

type Home struct {
	cwd              string
	basePath         string
	execPath         string
	lockSmith        interfaces.LockSmith
	age              *age.Age
	immutable_config immutable_config.Config
	debug            debug.Options
	dryRun           bool
	pid              int
}

func Make(
	// config immutable_config.Config,
	o Options,
) (s Home, err error) {
	s.age = &age.Age{}
	ui.TodoP3("add 'touched' which can get deleted / cleaned")
	if err = o.Validate(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if o.BasePath == "" {
		o.BasePath = os.Getenv(EnvDir)
	}

	if o.BasePath == "" {
		if o.BasePath, err = os.Getwd(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	s.dryRun = o.DryRun
	s.basePath = o.BasePath
	s.debug = o.Debug
	s.cwd = o.store_fs
	s.pid = os.Getpid()

	if !o.PermitNoZitDirectory {
		if ok := files.Exists(s.DirZit()); !ok {
			err = errors.Wrap(ErrNotInZitDir{})
			return
		}
	}

	s.lockSmith = file_lock.New(s.DirZit("Lock"))

	if s.execPath, err = os.Executable(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = os.Setenv(EnvDir, s.basePath); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = os.Setenv(EnvBin, s.execPath); err != nil {
		err = errors.Wrap(err)
		return
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

	return
}

func (a Home) SansAge() (b Home) {
	b = a
	b.age = nil
	return
}

func (a Home) SansCompression() (b Home) {
	b = a
	b.immutable_config.CompressionType = immutable_config.CompressionTypeNone
	return
}

func (s Home) GetConfig() immutable_config.Config {
	return s.immutable_config
}

func (s *Home) loadKonfigAngeboren() (err error) {
	var f *os.File

	if f, err = files.OpenExclusiveReadOnly(s.FileKonfigAngeboren()); err != nil {
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

func (s Home) GetLockSmith() interfaces.LockSmith {
	return s.lockSmith
}

func (s *Home) Age() *age.Age {
	return s.age
}

func (s Home) Cwd() string {
	return s.cwd
}

func (s Home) Executable() string {
	return s.execPath
}

func (s Home) AbsFromCwdOrSame(p string) (p1 string) {
	var err error
	p1, err = filepath.Abs(p)
	if err != nil {
		p1 = p
	}

	return
}

func (s Home) RelToCwdOrSame(p string) (p1 string) {
	var err error

	if p1, err = filepath.Rel(s.Cwd(), p); err != nil {
		p1 = p
	}

	return
}

func stringSliceJoin(s string, vs []string) []string {
	return append([]string{s}, vs...)
}

func (c Home) FileKonfigCompiled() string {
	return c.DirZit("KonfigCompiled")
}

func (c Home) FileSchlummernd() string {
	return c.DirZit("Schlummernd")
}

func (c Home) FileEtiketten() string {
	return c.DirZit("Etiketten")
}

func (c Home) FileKonfigAngeboren() string {
	return c.DirZit("KonfigAngeboren")
}

func (c Home) FileKonfigErworben() string {
	return c.DirZit("KonfigErworben")
}

func (c Home) FileKonfigToml() string {
	// var usr *user.User

	// if usr, err = user.Current(); err != nil {
	// 	err = errors.Error(err)
	// 	return
	// }

	// p = path.Join(
	// 	usr.HomeDir,
	// 	".config",
	// 	"zettelkasten",
	// 	"config.toml",
	// )

	return c.DirZit("Konfig")
}

func (s Home) Dir(p ...string) string {
	return filepath.Join(stringSliceJoin(s.basePath, p)...)
}

func (s Home) DirZit(p ...string) string {
	return s.Dir(stringSliceJoin(".zit", p)...)
}

func (s Home) FileAge() string {
	return s.DirZit("AgeIdentity")
}

func (s Home) DirVerzeichnisse(p ...string) string {
	return s.DirZit(append([]string{"Verzeichnisse"}, p...)...)
}

func (s Home) DirVerzeichnisseKasten(p ...string) string {
	return s.DirZit(append([]string{"Verzeichnisse", "Kasten"}, p...)...)
}

func (s Home) DirVerzeichnisseDurable(p ...string) string {
	return s.DirZit(append([]string{"VerzeichnisseDurable"}, p...)...)
}

func (s Home) DirObjekten(p ...string) string {
	return s.DirZit(append([]string{"Objekten"}, p...)...)
}

func (s Home) DirObjekten2(p ...string) string {
	return s.DirZit(append([]string{"Objekten2"}, p...)...)
}

func (s Home) DirVerlorenUndGefunden() string {
	return s.DirZit("Verloren+Gefunden")
}

func (s Home) MakeDir(d string) (err error) {
	if err = os.MkdirAll(d, os.ModeDir|0o755); err != nil {
		err = errors.Wrapf(err, "Dir: %q", d)
		return
	}

	return
}

func (s Home) DirVerzeichnisseObjekten() string {
	return s.DirVerzeichnisse("Objekten")
}

func (s Home) DirVerzeichnisseMetadata() string {
	return s.DirVerzeichnisseDurable("Metadatei")
}

func (s Home) DirVerzeichnisseMetadataObjectIdParent() string {
	return s.DirVerzeichnisseDurable("MetadateiKennungMutter")
}

func (s Home) DirVerzeichnisseVerweise() string {
	return s.DirVerzeichnisse("Verweise")
}

func (s Home) FileVerzeichnisseEtiketten() string {
	return s.DirVerzeichnisse("Etiketten")
}

func (s Home) FileVerzeichnisseObjectId() string {
	return s.DirVerzeichnisse("Kennung")
}

func (s Home) FileVerzeichnisseHinweis() string {
	return s.DirVerzeichnisse("Hinweis")
}

func (s Home) DirObjectId() string {
	return s.DirZit("Kennung")
}

func (s Home) ResetCache() (err error) {
	if err = files.SetAllowUserChangesRecursive(s.DirVerzeichnisse()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = os.RemoveAll(s.DirVerzeichnisse()); err != nil {
		err = errors.Wrapf(err, "failed to remove verzeichnisse dir")
		return
	}

	if err = s.MakeDir(s.DirVerzeichnisse()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.MakeDir(s.DirVerzeichnisseObjekten()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.MakeDir(s.DirVerzeichnisseMetadataObjectIdParent()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.MakeDir(s.DirVerzeichnisseVerweise()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (h Home) MakeCommonEnv() map[string]string {
	return map[string]string{
		"ZIT_BIN": h.Executable(),
		"ZIT_DIR": h.Dir(),
	}
}
