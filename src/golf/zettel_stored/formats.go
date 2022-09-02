package zettel_stored

import (
	"io"
)

type StoredFormat interface {
	ReadFrom(*Stored, io.Reader) (int64, error)
	WriteTo(Stored, io.Writer) (int64, error)
}
