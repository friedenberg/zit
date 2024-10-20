package immutable_config

import (
	"encoding/gob"
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
)

func init() {
	gob.RegisterName("KonfigLike", Config{})
}

type Config struct {
	StoreVersion                        storeVersion
	Recipients                          []string
	UseBestandsaufnahme                 bool // deprecated
	UseKonfigErworbenFile               bool // deprecated
	UseBestandsaufnahmeForVerzeichnisse bool // deprecated
	CompressionType                     CompressionType
	LockInternalFiles                   bool
}

func Default() Config {
	return Config{
		StoreVersion:      storeVersion(values.Int(7)),
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
