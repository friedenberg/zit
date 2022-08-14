package objekten

import "io"

type readCloserFactory interface {
	ReadCloserObjekten(string) (io.ReadCloser, error)
	ReadCloserVerzeichnisse(string) (io.ReadCloser, error)
}

type writeCloserFactory interface {
	WriteCloserObjekten(string) (io.WriteCloser, error)
	WriteCloserVerzeichnisse(string) (io.WriteCloser, error)
}

type ioFactory interface {
	readCloserFactory
	writeCloserFactory
}
