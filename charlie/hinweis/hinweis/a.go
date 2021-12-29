package hinweis

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/kennung"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/bravo/sha"
)

type (
	_Int     = kennung.Int
	_Kennung = kennung.Kennung
	_Sha     = sha.Sha
	_Id      = id.Id
)

var (
	_Errorf        = errors.Errorf
	_ReadAllString = open_file_guard.ReadAllString
	_Open          = open_file_guard.Open
	_Close         = open_file_guard.Close
	_IdPath        = id.Path
	_TempFile      = open_file_guard.TempFile
)
