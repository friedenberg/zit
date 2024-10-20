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

func (c directoryV0) FileCacheDormant() string {
	return c.DirZit("Schlummernd")
}

func (c directoryV0) FileTags() string {
	return c.DirZit("Etiketten")
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

type directoryV1 struct {
	basePath string
}

func (c directoryV1) GetDirectoryPaths() interfaces.DirectoryPaths {
	return c
}

func (c directoryV1) FileCacheDormant() string {
	return c.DirZit("Dormant")
}

func (c directoryV1) FileTags() string {
	return c.DirZit("Tags")
}

func (c directoryV1) FileConfigPermanent() string {
	return c.DirZit("ConfigPermanent")
}

func (c directoryV1) FileConfigMutable() string {
	return c.DirZit("ConfigMutable")
}

func (s directoryV1) Dir(p ...string) string {
	return filepath.Join(stringSliceJoin(s.basePath, p)...)
}

func (s directoryV1) DirZit(p ...string) string {
	return s.Dir(stringSliceJoin(".zit", p)...)
}

func (s directoryV1) FileAge() string {
	return s.DirZit("AgeIdentity")
}

func (s directoryV1) DirCache(p ...string) string {
	return s.DirZit(append([]string{"Cache"}, p...)...)
}

func (s directoryV1) DirCacheRepo(p ...string) string {
	return s.DirZit(append([]string{"Cache", "Repo"}, p...)...)
}

func (s directoryV1) DirObjects(p ...string) string {
	return s.DirZit(append([]string{"Objects"}, p...)...)
}

func (s directoryV1) DirLostAndFound() string {
	return s.DirZit("LostAndFound")
}

func (s directoryV1) DirCacheObjects() string {
	return s.DirCache("Objects")
}

func (s directoryV1) DirCacheObjectPointers() string {
	return s.DirCache("ObjectPointers")
}

func (s directoryV1) DirObjectId() string {
	return s.DirZit("ObjectId")
}

func (s directoryV1) FileCacheObjectId() string {
	return s.DirCache("ObjectId")
}
