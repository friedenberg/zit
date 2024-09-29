package interfaces

type Directory interface {
	Delete(string) error
	DirObjectId() string
	FileVerzeichnisseEtiketten() string
	FileVerzeichnisseObjectId() string
	FileVerzeichnisseHinweis() string
}

type CacheIOFactory interface {
	ReadCloserCache(string) (ShaReadCloser, error)
	WriteCloserCache(string) (ShaWriteCloser, error)
}
