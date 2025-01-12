package dir_layout

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
)

type Config struct {
	age               *age.Age
	compressionType   immutable_config.CompressionType
	lockInternalFiles bool
}

func MakeConfigFromImmutableBlobConfig(
	config immutable_config.BlobStoreConfig,
) Config {
	return Config{
		age:               config.GetAge(),
		compressionType:   config.GetCompressionType(),
		lockInternalFiles: config.GetLockInternalFiles(),
	}
}

func MakeConfig(
	age *age.Age,
	compressionType immutable_config.CompressionType,
	lockInternalFiles bool,
) Config {
	return Config{
		age:               age,
		compressionType:   compressionType,
		lockInternalFiles: lockInternalFiles,
	}
}

func (c Config) GetBlobStoreImmutableConfig() immutable_config.BlobStoreConfig {
	return c
}

func (c Config) GetAge() *age.Age {
	return &age.Age{}
}

func (c Config) GetCompressionType() immutable_config.CompressionType {
	return c.compressionType
}

func (c Config) GetLockInternalFiles() bool {
	return c.lockInternalFiles
}
