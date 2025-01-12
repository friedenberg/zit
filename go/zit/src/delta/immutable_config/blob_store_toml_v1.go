package immutable_config

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/age"
)

// TODO Split into repo and blob store configs
type BlobStoreTomlV1 struct {
	AgeIdentity       age.Age         `toml:"age-identity"`
	CompressionType   CompressionType `toml:"compression-type"`
	LockInternalFiles bool            `toml:"lock-internal-files"`
}

func (k *BlobStoreTomlV1) SetFlagSet(f *flag.FlagSet) {
	k.CompressionType.SetFlagSet(f)

	f.BoolVar(
		&k.LockInternalFiles,
		"lock-internal-files",
		k.LockInternalFiles,
		"",
	)

	f.Var(&k.AgeIdentity, "age-identity", "add an age identity")
}

func (k BlobStoreTomlV1) GetBlobStoreImmutableConfig() BlobStoreConfig {
	return k
}

func (k BlobStoreTomlV1) GetAge() *age.Age {
	return &age.Age{}
}

func (k BlobStoreTomlV1) GetCompressionType() CompressionType {
	return k.CompressionType
}

func (k BlobStoreTomlV1) GetLockInternalFiles() bool {
	return k.LockInternalFiles
}
