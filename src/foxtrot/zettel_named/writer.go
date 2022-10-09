package zettel_named

type Writer interface {
	WriteZettelNamed(*Zettel) (err error)
}
