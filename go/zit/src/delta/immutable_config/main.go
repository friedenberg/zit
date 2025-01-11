package immutable_config

import "code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"

type Latest = TomlV1

type Config interface {
	GetImmutableConfig() Config
	GetStoreVersion() interfaces.StoreVersion
	GetCompressionType() CompressionType
	GetLockInternalFiles() bool
}

func Default() TomlV1 {
	return TomlV1{
		StoreVersion:      CurrentStoreVersion,
		Recipients:        make([]string, 0),
		CompressionType:   CompressionTypeDefault,
		LockInternalFiles: true,
	}
}
