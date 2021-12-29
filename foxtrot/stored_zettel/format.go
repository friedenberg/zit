package stored_zettel

import "io"

type Format interface {
	ReadFrom(*Stored, io.Reader) (int64, error)
	WriteTo(Stored, io.Writer) (int64, error)
}
