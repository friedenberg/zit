package etikett

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/sha"
)

type (
	_Sha = sha.Sha
)

var (
	_Error           = errors.Error
	_Errorf          = errors.Errorf
	_MakeShaFromHash = sha.FromHash
)
