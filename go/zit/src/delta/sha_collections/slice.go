package sha_collections

import "github.com/friedenberg/zit/src/charlie/sha"

type Slice []sha.Sha

func MakeSlice(c int) Slice {
	return make([]sha.Sha, 0, c)
}

func (s *Slice) Append(sh ...sha.Sha) {
	*s = append(*s, sh...)
}
