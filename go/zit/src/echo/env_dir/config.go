package env_dir

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/config_immutable"
)

type Config struct {
	age               *age.Age
	compressionType   config_immutable.CompressionType
	lockInternalFiles bool
}

func MakeConfigFromImmutableBlobConfig(
	config config_immutable.BlobStoreConfig,
) Config {
	return MakeConfig(
		config.GetAgeEncryption(),
		config.GetCompressionType(),
		config.GetLockInternalFiles(),
	)
}

func MakeConfig(
	ag *age.Age,
	compressionType config_immutable.CompressionType,
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

func (c Config) GetBlobStoreImmutableConfig() config_immutable.BlobStoreConfig {
	return c
}

func (c Config) GetAgeEncryption() *age.Age {
	return c.age
}

func (c Config) GetCompressionType() config_immutable.CompressionType {
	return c.compressionType
}

func (c Config) GetLockInternalFiles() bool {
	return c.lockInternalFiles
}
