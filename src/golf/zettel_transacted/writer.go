package zettel_transacted

type Writer interface {
  WriteZettel(Zettel) (int, error)
}
