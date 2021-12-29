package akten

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/age"
	"github.com/friedenberg/zit/echo/sharded_store"
)

type (
	_Sha            = sha.Sha
	_Age            = age.Age
	_Store          = sharded_store.Store
	_StoreImmutable = sharded_store.StoreImmutable
	_Shard          = sharded_store.Shard
	_ShardImmutable = sharded_store.ShardImmutable
	_Entry          = sharded_store.Entry
	_ShardGeneric   = sharded_store.ShardGeneric
)

var (
	_NewStore          = sharded_store.NewStore
	_NewStoreImmutable = sharded_store.NewStoreImmutable
	_NewShard          = sharded_store.NewShard
	_NewShardImmutable = sharded_store.NewShardImmutable
	_Error             = errors.Error
	_Errorf            = errors.Errorf
	_ReadDirNames      = open_file_guard.ReadDirNames
)
