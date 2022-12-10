package sha

import "github.com/friedenberg/zit/src/bravo/sha"

const ShaNull = sha.ShaNull

type Sha = sha.Sha
type ReadCloser = sha.ReadCloser
type WriteCloser = sha.WriteCloser

var (
	FromHash        = sha.FromHash
	FromString      = sha.FromString
	MakeShaFromPath = sha.MakeShaFromPath
)
