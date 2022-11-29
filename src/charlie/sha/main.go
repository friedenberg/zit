package sha

import "github.com/friedenberg/zit/src/bravo/sha_core"

const ShaNull = sha_core.ShaNull

type Sha = sha_core.Sha
type ReadCloser = sha_core.ReadCloser
type WriteCloser = sha_core.WriteCloser

var (
	FromHash   = sha_core.FromHash
	FromString = sha_core.FromString
)
