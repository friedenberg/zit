package interfaces

type DirectoryPaths interface {
	Dir(p ...string) string
	DirCache(p ...string) string
	DirCacheDurable(p ...string) string
	DirCacheObjectPointers() string
	DirCacheRepo(p ...string) string
	DirObjectId() string
	DirObjects(p ...string) string
	DirVerlorenUndGefunden() string
	DirVerzeichnisseObjekten() string
	DirZit(p ...string) string
	FileAge() string
	FileCacheObjectId() string
	FileEtiketten() string
	FileKonfigAngeboren() string
	FileKonfigErworben() string
	FileSchlummernd() string
}

type Directory interface {
	DirectoryPaths
	Delete(string) error
}

type CacheIOFactory interface {
	ReadCloserCache(string) (ShaReadCloser, error)
	WriteCloserCache(string) (ShaWriteCloser, error)
}
