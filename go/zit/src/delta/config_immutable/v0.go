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

type v0Common struct {
	StoreVersion      StoreVersion
	Recipients        []string
	CompressionType   CompressionType
	LockInternalFiles bool
}

type V0Public struct {
  v0Common
}

type V0Private struct {
  v0Common
}

func (k *V0Public) SetFlagSet(f *flag.FlagSet) {
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
			k.Recipients = append(k.Recipients, value)
			return
		},
	)
}

func (k *V0Public) config() public { return public{} }
func (k *V0Private) config() private { return private{} }

func (k *V0Private) GetImmutableConfig() ConfigPrivate {
	return k
}

func (k *V0Private) GetImmutableConfigPublic() ConfigPublic {
	return &V0Public{
    v0Common: k.v0Common,
  }
}

func (k *V0Public) GetImmutableConfigPublic() ConfigPublic {
	return k
}

func (k *v0Common) GetBlobStoreConfigImmutable() interfaces.BlobStoreConfigImmutable {
	return k
}

func (k v0Common) GetStoreVersion() interfaces.StoreVersion {
	return k.StoreVersion
}

func (k v0Common) GetRepoType() repo_type.Type {
	return repo_type.TypeWorkingCopy
}

func (k v0Common) GetPrivateKey() ed25519.PrivateKey {
	panic(errors.Errorf("not supported"))
}

func (k v0Common) GetPublicKey() ed25519.PublicKey {
	panic(errors.Errorf("not supported"))
}

func (k v0Common) GetRepoId() ids.RepoId {
	return ids.RepoId{}
}

func (k *v0Common) GetAgeEncryption() *age.Age {
	return &age.Age{}
}

func (k *v0Common) GetBlobCompression() interfaces.BlobCompression {
	return &k.CompressionType
}

func (k *v0Common) GetBlobEncryption() interfaces.BlobEncryption {
	return k.GetAgeEncryption()
}

func (k v0Common) GetLockInternalFiles() bool {
	return k.LockInternalFiles
}
