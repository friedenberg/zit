package config_immutable

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/charlie/repo_signing"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type LatestPrivate = TomlV1Private

type (
	public  struct{}
	private struct{}
)

// TODO make it impossible for private configs to be returned fy
// GetImmutableConfigPublic
type configCommon interface {
	GetImmutableConfigPublic() ConfigPublic
	GetStoreVersion() interfaces.StoreVersion
	GetPublicKey() repo_signing.PublicKey
	GetRepoType() repo_type.Type
	GetRepoId() ids.RepoId
	GetBlobStoreConfigImmutable() interfaces.BlobStoreConfigImmutable
}

type ConfigPublic interface {
	config() public
	configCommon
}

type ConfigPrivate interface {
	configCommon
	config() private
	GetImmutableConfig() ConfigPrivate
	GetPrivateKey() repo_signing.PrivateKey
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
