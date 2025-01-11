package immutable_config

import "code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"

type Latest = V0

type Config interface {
	GetImmutableConfig() Config
	GetStoreVersion() interfaces.StoreVersion
	GetCompressionType() CompressionType
	GetLockInternalFiles() bool
}

func Default() V0 {
	return V0{
		StoreVersion:      CurrentStoreVersion,
		Recipients:        make([]string, 0),
		CompressionType:   CompressionTypeDefault,
		LockInternalFiles: true,
	}
}
