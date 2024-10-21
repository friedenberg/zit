package immutable_config

import (
	"encoding/gob"
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func init() {
	gob.RegisterName("KonfigLike", Config{})
}

type Config struct {
	StoreVersion                        StoreVersion
	Recipients                          []string
	UseBestandsaufnahme                 bool // deprecated
	UseKonfigErworbenFile               bool // deprecated
	UseBestandsaufnahmeForVerzeichnisse bool // deprecated
	CompressionType                     CompressionType
	LockInternalFiles                   bool
}

func Default() Config {
	return Config{
		StoreVersion:      CurrentStoreVersion,
		Recipients:        make([]string, 0),
		CompressionType:   CompressionTypeDefault,
		LockInternalFiles: true,
	}
}

func (k Config) GetStoreVersion() interfaces.StoreVersion {
	return k.StoreVersion
}

func (k *Config) AddToFlagSet(f *flag.FlagSet) {
	k.CompressionType.AddToFlagSet(f)

	f.BoolVar(
		&k.LockInternalFiles,
		"lock-internal-files",
		k.LockInternalFiles,
		"",
	)

	f.Func(
		"recipient",
		"age recipients",
		func(value string) (err error) {
			// TODO-P2 validate age recipient
			k.Recipients = append(k.Recipients, value)
			return
		},
	)
}
