package fs_home

import (
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type directoryV1 struct {
	XDG
}

func (c *directoryV1) init(xdg XDG) (err error) {
	c.XDG = xdg
	return
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
	return filepath.Join(stringSliceJoin(s.Data, p)...)
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
