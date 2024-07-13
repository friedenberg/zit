package interfaces

type VerzeichnisseFactory interface {
	ReadCloserVerzeichnisse(string) (ShaReadCloser, error)
	WriteCloserVerzeichnisse(string) (ShaWriteCloser, error)
}
