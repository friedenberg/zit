package script_value

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/open_file_guard"
)

var (
	_Open  = open_file_guard.Open
	_Close = open_file_guard.Close
	_Error = errors.Error
)
