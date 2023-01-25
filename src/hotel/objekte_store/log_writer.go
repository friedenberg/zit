package objekte_store

import "github.com/friedenberg/zit/src/charlie/collections"

type LogWriter[
	T any,
] struct {
	New, Updated, Unchanged, Archived collections.WriterFunc[T]
}
