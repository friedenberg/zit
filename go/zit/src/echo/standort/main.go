package standort

import (
	"encoding/gob"
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/angeboren"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_lock"
)

type Getter interface {
	GetStandort() Standort
}

type Standort struct {
	cwd       string
	basePath  string
	execPath  string
	lockSmith schnittstellen.LockSmith
	age       *age.Age
	angeboren angeboren.Konfig
	debug     debug.Options
	dryRun    bool
	pid       int
}

func Make(
	o Options,
) (s Standort, err error) {
	s.age = &age.Age{}
	errors.TodoP3("add 'touched' which can get deleted / cleaned")
	if err = o.Validate(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if o.BasePath == "" {
		o.BasePath = os.Getenv("DIR_ZIT")
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
	s.cwd = o.cwd
	s.pid = os.Getpid()

	if ok := files.Exists(s.DirZit()); !ok {
		err = errors.Wrap(ErrNotInZitDir{})
		return
	}

	s.lockSmith = file_lock.New(s.DirZit("Lock"))

	if s.execPath, err = os.Executable(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = os.Setenv("DIR_ZIT", s.basePath); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = os.Setenv("BIN_ZIT", s.execPath); err != nil {
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

func (a Standort) SansAge() (b Standort) {
	b = a
	b.age = nil
	return
}

func (a Standort) SansCompression() (b Standort) {
	b = a
	b.angeboren.CompressionType = angeboren.CompressionTypeNone
	return
}

func (s Standort) GetKonfig() angeboren.Konfig {
	return s.angeboren
}

func (s *Standort) loadKonfigAngeboren() (err error) {
	var f *os.File

	if f, err = files.OpenExclusiveReadOnly(s.FileKonfigAngeboren()); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.Deferred(&err, f.Close)

	dec := gob.NewDecoder(f)

	if err = dec.Decode(&s.angeboren); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Standort) GetLockSmith() schnittstellen.LockSmith {
	return s.lockSmith
}

func (s *Standort) Age() *age.Age {
	return s.age
}

func (s Standort) Cwd() string {
	return s.cwd
}

func (s Standort) Executable() string {
	return s.execPath
}

func (s Standort) RelToCwdOrSame(p string) (p1 string) {
	var err error
	p1, err = filepath.Rel(s.Cwd(), p)
	if err != nil {
		p1 = p
	}

	return
}

func stringSliceJoin(s string, vs []string) []string {
	return append([]string{s}, vs...)
}

func (c Standort) FileKonfigCompiled() string {
	return c.DirZit("KonfigCompiled")
}

func (c Standort) FileSchlummernd() string {
	return c.DirZit("Schlummernd")
}

func (c Standort) FileEtiketten() string {
	return c.DirZit("Etiketten")
}

func (c Standort) FileKonfigAngeboren() string {
	return c.DirZit("KonfigAngeboren")
}

func (c Standort) FileKonfigErworben() string {
	return c.DirZit("KonfigErworben")
}

func (c Standort) FileKonfigToml() string {
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

func (s Standort) Dir(p ...string) string {
	return filepath.Join(stringSliceJoin(s.basePath, p)...)
}

func (s Standort) DirZit(p ...string) string {
	return s.Dir(stringSliceJoin(".zit", p)...)
}

func (s Standort) FileAge() string {
	return s.DirZit("AgeIdentity")
}

func (s Standort) DirVerzeichnisse(p ...string) string {
	return s.DirZit(append([]string{"Verzeichnisse"}, p...)...)
}

func (s Standort) DirVerzeichnisseDurable(p ...string) string {
	return s.DirZit(append([]string{"VerzeichnisseDurable"}, p...)...)
}

func (s Standort) DirObjekten(p ...string) string {
	return s.DirZit(append([]string{"Objekten"}, p...)...)
}

func (s Standort) DirObjekten2(p ...string) string {
	return s.DirZit(append([]string{"Objekten2"}, p...)...)
}

func (s Standort) DirVerlorenUndGefunden() string {
	return s.DirZit("Verloren+Gefunden")
}

func (s Standort) MakeDir(d string) (err error) {
	if err = os.MkdirAll(d, os.ModeDir|0o755); err != nil {
		err = errors.Wrapf(err, "Dir: %q", d)
		return
	}

	return
}

func (s Standort) DirVerzeichnisseObjekten() string {
	return s.DirVerzeichnisse("Objekten")
}

func (s Standort) DirVerzeichnisseMetadatei() string {
	return s.DirVerzeichnisseDurable("Metadatei")
}

func (s Standort) DirVerzeichnisseMetadateiKennungMutter() string {
	return s.DirVerzeichnisseDurable("MetadateiKennungMutter")
}

func (s Standort) DirVerzeichnisseVerweise() string {
	return s.DirVerzeichnisse("Verweise")
}

func (s Standort) FileVerzeichnisseEtiketten() string {
	return s.DirVerzeichnisse("Etiketten")
}

func (s Standort) FileVerzeichnisseKennung() string {
	return s.DirVerzeichnisse("Kennung")
}

func (s Standort) FileVerzeichnisseHinweis() string {
	return s.DirVerzeichnisse("Hinweis")
}

func (s Standort) DirKennung() string {
	return s.DirZit("Kennung")
}

func (s Standort) ResetVerzeichnisse() (err error) {
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

	if err = s.MakeDir(s.DirVerzeichnisseMetadateiKennungMutter()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.MakeDir(s.DirVerzeichnisseVerweise()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
