package stored_zettel

import "io"

type StoredFormat interface {
	ReadFrom(*Stored, io.Reader) (int64, error)
	WriteTo(Stored, io.Writer) (int64, error)
}

type ExternalFormat interface {
	ReadExternalZettelFrom(*External, io.Reader) (int64, error)
	WriteExternalZettelTo(External, io.Writer) (int64, error)
}
