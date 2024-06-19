package schnittstellen

type Standort interface {
	Delete(string) error
	DirKennung() string
	FileVerzeichnisseEtiketten() string
	FileVerzeichnisseKennung() string
	FileVerzeichnisseHinweis() string
}
