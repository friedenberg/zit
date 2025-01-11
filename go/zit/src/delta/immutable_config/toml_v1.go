package immutable_config

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
)

// TODO Split into repo and blob store configs
type TomlV1 struct {
	StoreVersion      StoreVersion    `toml:"store-version"`
	RepoType          repo_type.Type  `toml:"repo-type"`
	Recipients        []string        `toml:"recipients"`
	CompressionType   CompressionType `toml:"compression-type"`
	LockInternalFiles bool            `toml:"lock-internal-files"`
}

func (k *TomlV1) SetFlagSet(f *flag.FlagSet) {
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

func (k TomlV1) GetImmutableConfig() Config {
	return k
}

func (k TomlV1) GetStoreVersion() interfaces.StoreVersion {
	return k.StoreVersion
}

func (k TomlV1) GetCompressionType() CompressionType {
	return k.CompressionType
}

func (k TomlV1) GetLockInternalFiles() bool {
	return k.LockInternalFiles
}
