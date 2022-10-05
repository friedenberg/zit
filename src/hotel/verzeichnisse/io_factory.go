package verzeichnisse

import "io"

type readCloserFactory interface {
	ReadCloserVerzeichnisse(string) (io.ReadCloser, error)
}

type writeCloserFactory interface {
	WriteCloserVerzeichnisse(string) (io.WriteCloser, error)
}

type ioFactory interface {
	readCloserFactory
	writeCloserFactory
}
