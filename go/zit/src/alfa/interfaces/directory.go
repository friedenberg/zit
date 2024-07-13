package interfaces

type Directory interface {
	Delete(string) error
	DirKennung() string
	FileVerzeichnisseEtiketten() string
	FileVerzeichnisseKennung() string
	FileVerzeichnisseHinweis() string
}

type CacheIOFactory interface {
	ReadCloserCache(string) (ShaReadCloser, error)
	WriteCloserCache(string) (ShaWriteCloser, error)
}
