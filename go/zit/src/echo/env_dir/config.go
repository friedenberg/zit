package env_dir

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
	return MakeConfig(
		config.GetAgeEncryption(),
		config.GetCompressionType(),
		config.GetLockInternalFiles(),
	)
}

func MakeConfig(
	ag *age.Age,
	compressionType immutable_config.CompressionType,
	lockInternalFiles bool,
) Config {
	if ag == nil {
		ag = &age.Age{}
	}

	return Config{
		age:               ag,
		compressionType:   compressionType,
		lockInternalFiles: lockInternalFiles,
	}
}

func (c Config) GetBlobStoreImmutableConfig() immutable_config.BlobStoreConfig {
	return c
}

func (c Config) GetAgeEncryption() *age.Age {
	return c.age
}

func (c Config) GetCompressionType() immutable_config.CompressionType {
	return c.compressionType
}

func (c Config) GetLockInternalFiles() bool {
	return c.lockInternalFiles
}
