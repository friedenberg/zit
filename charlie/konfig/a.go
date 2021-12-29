package konfig

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/node_type"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/alfa/toml"
	"github.com/friedenberg/zit/bravo/open_file_guard"
)

type (
	_Type   = node_type.Type
	_Logger = stdprinter.Logger
)

var (
	_Open          = open_file_guard.Open
	_Close         = open_file_guard.Close
	_Errorf        = errors.Errorf
	_Error         = errors.Error
	_TomlUnmarshal = toml.Unmarshal
)
