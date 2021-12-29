package umwelt

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/age"
	"github.com/friedenberg/zit/charlie/file_lock"
	"github.com/friedenberg/zit/charlie/konfig"
)

type (
	_Konfig   = konfig.Konfig
	_Logger   = stdprinter.Logger
	_Age      = age.Age
	_FileLock = file_lock.Lock
)

var (
	_Close       = open_file_guard.Close
	_Open        = open_file_guard.Open
	_AgeMake     = age.Make
	_Error       = errors.Error
	_FileLockNew = file_lock.New
)
