package zettel_transacted

type Writer interface {
  WriteZettelTransacted(Zettel) (error)
}
