package sharded_store

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/age"
	"github.com/friedenberg/zit/delta/objekte"
)

type (
	_Id            = id.Id
	_Age           = age.Age
	_ObjekteReader = objekte.Reader
	_ObjekteWriter = objekte.Writer
)

var (
	_Error            = errors.Error
	_OpenFile         = open_file_guard.OpenFile
	_Close            = open_file_guard.Close
	_ReadDirNames     = open_file_guard.ReadDirNames
	_ReadAllString    = open_file_guard.ReadAllString
	_NewObjekteReader = objekte.NewReader
	_NewObjekteWriter = objekte.NewWriter
	_TempFile         = open_file_guard.TempFile
	_PanicIfError     = stdprinter.PanicIfError
	_Errorf           = errors.Errorf
)
