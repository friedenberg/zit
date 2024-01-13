package schnittstellen

type Standort interface {
	DirKennung() string
	FileVerzeichnisseEtiketten() string
	FileVerzeichnisseKennung() string
	FileVerzeichnisseHinweis() string
}
