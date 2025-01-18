package config_immutable

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
)

// TODO Split into repo and blob store configs
type BlobStoreTomlV1 struct {
	AgeEncryption     age.Age         `toml:"age-encryption,omitempty"`
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

	f.Var(&k.AgeEncryption, "age-identity", "add an age identity")
}

func (k *BlobStoreTomlV1) GetBlobStoreImmutableConfig() interfaces.BlobStoreConfig {
	return k
}

func (k *BlobStoreTomlV1) GetBlobCompression() interfaces.BlobCompression {
	return &k.CompressionType
}

func (k *BlobStoreTomlV1) GetBlobEncryption() interfaces.BlobEncryption {
	return &k.AgeEncryption
}

func (k *BlobStoreTomlV1) GetAgeEncryption() *age.Age {
	return &k.AgeEncryption
}

func (k *BlobStoreTomlV1) GetLockInternalFiles() bool {
	return k.LockInternalFiles
}
