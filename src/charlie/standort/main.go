package standort

import (
	"os"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
)

type Standort struct {
	cwd      string
	basePath string
	execPath string
}

func Make(o Options) (s Standort, err error) {
	errors.TodoP3("add 'touched' which can get deleted / cleaned")
	if err = o.Validate(); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.basePath = o.BasePath
	s.cwd = o.cwd

	if ok := files.Exists(s.DirZit()); !ok {
		err = errors.Wrap(ErrNotInZitDir{})
		return
	}

	if s.execPath, err = os.Executable(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
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

func (s Standort) DirObjekten(p ...string) string {
	return s.DirZit(append([]string{"Objekten"}, p...)...)
}

// TODO spelling?
func (s Standort) DirVerlorenUndGefunden() string {
	return s.DirZit("Verloren+Gefunden")
}

func (s Standort) FileVerzeichnisseZettelenSchwanzen() string {
	return s.DirVerzeichnisse("ZettelenSchwanzen")
}

func (s Standort) DirVerzeichnisseZettelenNeue() string {
	return s.DirVerzeichnisse("ZettelenNeue")
}

func (s Standort) DirVerzeichnisseZettelenNeueSchwanzen() string {
	return s.DirVerzeichnisse("ZettelenNeueSchwanzen")
}

func (s Standort) DirVerzeichnisseAkten() string {
	return s.DirVerzeichnisse("Akten")
}

func (s Standort) FileVerzeichnisseZettelen() string {
	return s.DirVerzeichnisse("Zettelen")
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
	if err = os.RemoveAll(s.DirVerzeichnisse()); err != nil {
		err = errors.Wrapf(err, "failed to remove verzeichnisse dir")
		return
	}

	if err = os.MkdirAll(s.DirVerzeichnisse(), os.ModeDir|0o755); err != nil {
		err = errors.Wrapf(err, "failed to make verzeichnisse dir")
		return
	}

	if err = os.MkdirAll(s.DirVerzeichnisseZettelenNeue(), os.ModeDir|0o755); err != nil {
		err = errors.Wrapf(err, "failed to make verzeichnisse dir")
		return
	}

	if err = os.MkdirAll(s.DirVerzeichnisseZettelenNeueSchwanzen(), os.ModeDir|0o755); err != nil {
		err = errors.Wrapf(err, "failed to make verzeichnisse dir")
		return
	}

	return
}
