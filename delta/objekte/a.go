package objekte

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/node_type"
	"github.com/friedenberg/zit/bravo/files"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/age"
)

type (
	_Sha  = sha.Sha
	_Type = node_type.Type
	_Age  = age.Age
)

var (
	_Close                       = open_file_guard.Close
	_Errorf                      = errors.Errorf
	_Error                       = errors.Error
	_MakeShaFromHash             = sha.FromHash
	_Open                        = open_file_guard.Open
	_TempFile                    = open_file_guard.TempFile
	_IdMakeDirIfNecessary        = id.MakeDirIfNecessary
	_FilesExists                 = files.Exists
	_FilesSetDisallowUserChanges = files.SetDisallowUserChanges
)
