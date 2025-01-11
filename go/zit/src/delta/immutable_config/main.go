package immutable_config

import (
	"encoding/gob"
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
)

func init() {
	gob.RegisterName("KonfigLike", Config{})
}

// TODO Split into repo and blob store configs
// TODO make toml-compatible
type Config struct {
	StoreVersion      StoreVersion    `toml:"store-version"`
	RepoType          repo_type.Type  `toml:"repo-type"`
	Recipients        []string        `toml:"recipients"`
	CompressionType   CompressionType `toml:"compression-type"`
	LockInternalFiles bool            `toml:"lock-internal-files"`
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

func (k *Config) SetFlagSet(f *flag.FlagSet) {
	k.CompressionType.SetFlagSet(f)

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
