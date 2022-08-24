package collections

import (
	"github.com/friedenberg/zit/src/golf/stored_zettel"
)

type SliceTransacted []stored_zettel.Transacted

func MakeSliceTransacted() SliceTransacted {
	return make(SliceTransacted, 0)
}

func (s SliceTransacted) Len() int {
  return len(s)
}
