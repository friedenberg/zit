package akten

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/age"
	"github.com/friedenberg/zit/echo/sharded_store"
)

type (
	_Sha          = sha.Sha
	_Age          = age.Age
	_Store        = sharded_store.Store
	_Shard        = sharded_store.Shard
	_Entry        = sharded_store.Entry
	_ShardGeneric = sharded_store.ShardGeneric
)

var (
	_NewStore     = sharded_store.NewStore
	_NewShard     = sharded_store.NewShard
	_Error        = errors.Error
	_Errorf       = errors.Errorf
	_ReadDirNames = open_file_guard.ReadDirNames
)
