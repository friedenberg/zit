package dir_layout

import (
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/xdg"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
)

type directoryPaths interface {
	interfaces.DirectoryPaths
	init(sv immutable_config.StoreVersion, xdg xdg.XDG) error
}

type directoryV0 struct {
	sv       immutable_config.StoreVersion
	basePath string
}

func (c *directoryV0) init(
	sv immutable_config.StoreVersion,
	xdg xdg.XDG,
) (err error) {
	c.sv = sv
	return
}

func (c directoryV0) GetDirectoryPaths() interfaces.DirectoryPaths {
	return c
}

func (c directoryV0) FileCacheDormant() string {
	return c.DirZit("Schlummernd")
}

func (c directoryV0) FileTags() string {
	return c.DirZit("Etiketten")
}

func (c directoryV0) FileLock() string {
	return c.DirZit("Lock")
}

func (c directoryV0) FileConfigPermanent() string {
	return c.DirZit("KonfigAngeboren")
}

func (c directoryV0) FileConfigMutable() string {
	return c.DirZit("KonfigErworben")
}

func (s directoryV0) Dir(p ...string) string {
	return filepath.Join(stringSliceJoin(s.basePath, p)...)
}

func (s directoryV0) DirZit(p ...string) string {
	return s.Dir(stringSliceJoin(".zit", p)...)
}

func (s directoryV0) DirObjectGenre(
	g1 interfaces.GenreGetter,
) (p string, err error) {
	g := g1.GetGenre()

	if g == genres.None {
		err = genres.MakeErrUnsupportedGenre(g)
		return
	}

	p = s.DirObjects(g.GetGenreStringPlural(s.sv))

	return
}

func (s directoryV0) FileAge() string {
	return s.DirZit("AgeIdentity")
}

func (s directoryV0) DirCache(p ...string) string {
	return s.DirZit(append([]string{"Verzeichnisse"}, p...)...)
}

func (s directoryV0) DirCacheRepo(p ...string) string {
	return s.DirZit(append([]string{"Verzeichnisse", "Kasten"}, p...)...)
}

func (s directoryV0) DirCacheDurable(p ...string) string {
	return s.DirZit(append([]string{"VerzeichnisseDurable"}, p...)...)
}

func (s directoryV0) DirObjects(p ...string) string {
	return s.DirZit(append([]string{"Objekten2"}, p...)...)
}

func (s directoryV0) DirLostAndFound() string {
	return s.DirZit("Verloren+Gefunden")
}

func (s directoryV0) DirCacheObjects() string {
	return s.DirCache("Objekten")
}

func (s directoryV0) DirCacheObjectPointers() string {
	return s.DirCache("Verweise")
}

func (s directoryV0) DirObjectId() string {
	return s.DirZit("Kennung")
}

func (s directoryV0) FileCacheObjectId() string {
	return s.DirCache("Kennung")
}
