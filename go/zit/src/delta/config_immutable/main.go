package config_immutable

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
)

type Latest = TomlV1

type Config interface {
	GetImmutableConfig() Config
	GetStoreVersion() interfaces.StoreVersion
	GetRepoType() repo_type.Type
	GetBlobStoreConfigImmutable() interfaces.BlobStoreConfigImmutable
}

type BlobStoreConfig interface {
	interfaces.BlobStoreConfigImmutable
	GetCompressionType() CompressionType
	GetAgeEncryption() *age.Age
	GetLockInternalFiles() bool
}

func Default() *TomlV1 {
	return &TomlV1{
		StoreVersion: CurrentStoreVersion,
		RepoType:     repo_type.TypeWorkingCopy,
		BlobStore: BlobStoreTomlV1{
			CompressionType:   CompressionTypeDefault,
			LockInternalFiles: true,
		},
	}
}
