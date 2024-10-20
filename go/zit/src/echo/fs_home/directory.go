package fs_home

import (
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type directoryV0 struct {
	basePath string
}

func (c directoryV0) GetDirectoryPaths() interfaces.DirectoryPaths {
  return c
}

func (c directoryV0) FileSchlummernd() string {
	return c.DirZit("Schlummernd")
}

func (c directoryV0) FileEtiketten() string {
	return c.DirZit("Etiketten")
}

func (c directoryV0) FileKonfigAngeboren() string {
	return c.DirZit("KonfigAngeboren")
}

func (c directoryV0) FileKonfigErworben() string {
	return c.DirZit("KonfigErworben")
}

func (s directoryV0) Dir(p ...string) string {
	return filepath.Join(stringSliceJoin(s.basePath, p)...)
}

func (s directoryV0) DirZit(p ...string) string {
	return s.Dir(stringSliceJoin(".zit", p)...)
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

func (s directoryV0) DirObjektenOld(p ...string) string {
	return s.DirZit(append([]string{"Objekten"}, p...)...)
}

func (s directoryV0) DirObjects(p ...string) string {
	return s.DirZit(append([]string{"Objekten2"}, p...)...)
}

func (s directoryV0) DirVerlorenUndGefunden() string {
	return s.DirZit("Verloren+Gefunden")
}

func (s directoryV0) DirVerzeichnisseObjekten() string {
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
