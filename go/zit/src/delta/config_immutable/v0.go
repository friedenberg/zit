package config_immutable

import (
	"crypto/ed25519"
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
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

func (k *V0) GetImmutableConfig() Config {
	return k
}

func (k *V0) GetBlobStoreConfigImmutable() interfaces.BlobStoreConfigImmutable {
	return k
}

func (k V0) GetStoreVersion() interfaces.StoreVersion {
	return k.StoreVersion
}

func (k V0) GetRepoType() repo_type.Type {
	return repo_type.TypeWorkingCopy
}

func (k V0) GetPrivateKey() ed25519.PrivateKey {
	panic(errors.Errorf("not supported"))
}

func (k V0) GetRepoId() ids.RepoId {
	return ids.RepoId{}
}

func (k *V0) GetAgeEncryption() *age.Age {
	return &age.Age{}
}

func (k *V0) GetBlobCompression() interfaces.BlobCompression {
	return &k.CompressionType
}

func (k *V0) GetBlobEncryption() interfaces.BlobEncryption {
	return k.GetAgeEncryption()
}

func (k V0) GetLockInternalFiles() bool {
	return k.LockInternalFiles
}
