package immutable_config

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
)

type Latest = TomlV1

type Config interface {
	GetImmutableConfig() Config
	GetStoreVersion() interfaces.StoreVersion
	GetBlobStoreImmutableConfig() BlobStoreConfig
}

type BlobStoreConfig interface {
	GetBlobStoreImmutableConfig() BlobStoreConfig
	GetAge() age.Age
	GetCompressionType() CompressionType
	GetLockInternalFiles() bool
}

func Default() TomlV1 {
	return TomlV1{
		StoreVersion: CurrentStoreVersion,
		RepoType:     repo_type.TypeReadWrite,
		BlobStore: BlobStoreTomlV1{
			Recipients:        make([]string, 0),
			CompressionType:   CompressionTypeDefault,
			LockInternalFiles: true,
		},
	}
}
