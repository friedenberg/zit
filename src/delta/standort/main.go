package standort

import (
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
)

type Standort struct {
	cwd      string
	basePath string
}

func Make(o Options) (s Standort, err error) {
	if err = o.Validate(); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.basePath = o.BasePath
	s.cwd = o.cwd

	if ok := files.Exists(s.DirZit()); !ok {
		err = ErrNotInZitDir{}
		return
	}

	return
}

func (s Standort) Cwd() string {
	return s.cwd
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
func (s Standort) DirObjektenKennungen() string {
	return s.DirObjekten("Kennungen")
}

func (s Standort) DirObjektenZettelen() string {
	return s.DirObjekten("Zettelen")
}

func (s Standort) DirObjektenTransaktion() string {
	return s.DirObjekten("Transaktion")
}

func (s Standort) DirObjektenAkten() string {
	return s.DirObjekten("Akten")
}

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
