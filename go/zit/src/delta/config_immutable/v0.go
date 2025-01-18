package config_immutable

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
)

// TODO Split into repo and blob store configs
type V0 struct {
	StoreVersion      StoreVersion
	Recipients        []string
	CompressionType   CompressionType
	LockInternalFiles bool
}

func (k *V0) SetFlagSet(f *flag.FlagSet) {
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

func (k V0) GetImmutableConfig() Config {
	return k
}

func (k V0) GetBlobStoreImmutableConfig() BlobStoreConfig {
	return k
}

func (k V0) GetStoreVersion() interfaces.StoreVersion {
	return k.StoreVersion
}

func (k V0) GetRepoType() repo_type.Type {
	return repo_type.TypeWorkingCopy
}

func (k V0) GetAgeEncryption() *age.Age {
	return &age.Age{}
}

func (k V0) GetCompressionType() CompressionType {
	return k.CompressionType
}

func (k V0) GetLockInternalFiles() bool {
	return k.LockInternalFiles
}
