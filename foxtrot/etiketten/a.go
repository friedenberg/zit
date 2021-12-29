package etiketten

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/age"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/charlie/konfig"
	"github.com/friedenberg/zit/delta/objekte"
	"github.com/friedenberg/zit/echo/sharded_store"
)

type (
	_Age           = age.Age
	_Etikett       = etikett.Etikett
	_Sha           = sha.Sha
	_ObjekteReader = objekte.Reader
	_Shard         = sharded_store.Shard
	_Entry         = sharded_store.Entry
	_Store         = sharded_store.Store
	_ShardLine     = sharded_store.ShardLine
	_ShardGeneric  = sharded_store.ShardGeneric
	_Konfig        = konfig.Konfig
)

var (
	_Open                = open_file_guard.Open
	_OpenFile            = open_file_guard.OpenFile
	_Close               = open_file_guard.Close
	_ObjekteDecodeBase64 = objekte.DecodeBase64
	_ObjekteEncodeBase64 = objekte.EncodeBase64
	_Errorf              = errors.Errorf
	_Error               = errors.Error
	_Err                 = stdprinter.Err
	_NewShard            = sharded_store.NewShard
	_NewStore            = sharded_store.NewStore
)
