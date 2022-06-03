package umwelt

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/files"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/age"
	"github.com/friedenberg/zit/charlie/file_lock"
)

var (
	_Close        = open_file_guard.Close
	_Open         = open_file_guard.Open
	_AgeMake      = age.Make
	_AgeMakeEmpty = age.MakeEmpty
	_Error        = errors.Error
	_FileLockNew  = file_lock.New
	_FilesExist   = files.Exists
)
