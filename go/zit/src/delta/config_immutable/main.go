package config_immutable

import (
	"crypto/ed25519"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type LatestPrivate = TomlV1Private

// TODO make it impossible for private configs to be returned fy
// GetImmutableConfigPublic
type ConfigPublic interface {
	GetImmutableConfigPublic() ConfigPublic
	GetStoreVersion() interfaces.StoreVersion
	GetPublicKey() ed25519.PublicKey
	GetRepoType() repo_type.Type
	GetRepoId() ids.RepoId
	GetBlobStoreConfigImmutable() interfaces.BlobStoreConfigImmutable
}

type ConfigPrivate interface {
	ConfigPublic
	GetImmutableConfig() ConfigPrivate
	GetPrivateKey() ed25519.PrivateKey
}

type BlobStoreConfig interface {
	interfaces.BlobStoreConfigImmutable
	GetCompressionType() CompressionType
	GetAgeEncryption() *age.Age
	GetLockInternalFiles() bool
}

func Default() *TomlV1Private {
	return &TomlV1Private{
		TomlV1Common: TomlV1Common{
			StoreVersion: CurrentStoreVersion,
			RepoType:     repo_type.TypeWorkingCopy,
			BlobStore: BlobStoreTomlV1{
				CompressionType:   CompressionTypeDefault,
				LockInternalFiles: true,
			},
		},
	}
}
