package objekte_store

import "github.com/friedenberg/zit/src/charlie/collections"

type LogWriter[
	T any,
] struct {
	New, Updated, Unchanged, Archived collections.WriterFunc[T]
}

func (l LogWriter[T]) NewOrUpdated(err error) collections.WriterFunc[T] {
  if IsNotFound(err) {
    return l.New
  } else {
    return l.Updated
  }
}
