package file_lock

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/files"
	"github.com/friedenberg/zit/bravo/open_file_guard"
)

var (
	_FileExists = files.Exists
	_Error      = errors.Error
	_Errorf     = errors.Errorf
	_OpenFile   = open_file_guard.OpenFile
	_Close      = open_file_guard.Close
)
