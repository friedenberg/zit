package interfaces

type DirectoryPaths interface {
	Dir(p ...string) string
	DirCache(p ...string) string
	DirCacheObjectPointers() string
	DirCacheObjects() string
	DirCacheRepo(p ...string) string
	DirLostAndFound() string
	DirObjectId() string
	DirObjects(p ...string) string
	DirZit(p ...string) string
	FileAge() string
	FileCacheDormant() string
	FileCacheObjectId() string
	FileConfigMutable() string
	FileConfigPermanent() string
	FileTags() string
}

type Directory interface {
	DirectoryPaths
	Delete(string) error
}

type CacheIOFactory interface {
	ReadCloserCache(string) (ShaReadCloser, error)
	WriteCloserCache(string) (ShaWriteCloser, error)
}
