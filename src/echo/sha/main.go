package sha

import (
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/gattung"
)

const ShaNull = sha.ShaNull

type Sha = sha.Sha
type ShaLike = gattung.ShaLike
type ReadCloser = sha.ReadCloser
type WriteCloser = sha.WriteCloser

var (
	FromHash          = sha.FromHash
	FromString        = sha.FromString
	MakeHashWriter    = sha.MakeHashWriter
	MakeReadCloser    = sha.MakeReadCloser
	MakeReadCloserTee = sha.MakeReadCloserTee
	MakeShaFromPath   = sha.MakeShaFromPath
)
