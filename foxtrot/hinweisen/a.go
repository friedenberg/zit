package hinweisen

import (
	"github.com/friedenberg/zit/alfa/kennung"
	"github.com/friedenberg/zit/alfa/node_type"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/age"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/delta/objekte"
	"github.com/friedenberg/zit/echo/sharded_store"
)

const (
	_TypeAkte = node_type.TypeAkte
)

type (
	_Age          = age.Age
	_Sha          = sha.Sha
	_Int          = kennung.Int
	_Hinweis      = hinweis.Hinweis
	_Store        = sharded_store.Store
	_Shard        = sharded_store.Shard
	_ShardGeneric = sharded_store.ShardGeneric
	_Entry        = sharded_store.Entry
)

var (
	_ObjekteWriteAndMove = objekte.WriteAndMove
	_ObjekteEncodeBase64 = objekte.EncodeBase64
	_ObjekteDecodeBase64 = objekte.DecodeBase64
	_Open                = open_file_guard.Open
	_OpenFile            = open_file_guard.OpenFile
	_Close               = open_file_guard.Close
	_Err                 = stdprinter.Err
	_MakeBlindHinweis    = hinweis.MakeBlindHinweis
	_ReadAllString       = open_file_guard.ReadAllString
	_TempFile            = open_file_guard.TempFile
	_NewStore            = sharded_store.NewStore
	_NewShard            = sharded_store.NewShard
)
