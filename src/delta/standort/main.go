package standort

import (
	"os"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/konfig"
)

type Standort struct {
	cwd      string
	basePath string
}

func Make(k konfig.Konfig) (s Standort, err error) {
	if s.basePath, err = k.DirZit(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.cwd, err = os.Getwd(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if ok := files.Exists(s.DirZit()); !ok {
		err = errors.Normalf("not in a zit directory")
		return
	}

	return
}

func (s Standort) Cwd() string {
	return s.cwd
}

func stringSliceJoin(s string, vs []string) []string {
	return append([]string{s}, vs...)
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
