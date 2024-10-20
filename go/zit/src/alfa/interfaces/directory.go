package interfaces

type Directory interface {
	Delete(string) error
	DirObjectId() string
	FileCacheObjectId() string
}

type CacheIOFactory interface {
	ReadCloserCache(string) (ShaReadCloser, error)
	WriteCloserCache(string) (ShaWriteCloser, error)
}
