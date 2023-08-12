package angeboren

import (
	"encoding/gob"
	"flag"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
)

func init() {
	gob.RegisterName("KonfigLike", Konfig{})
}

type Konfig struct {
	StoreVersion          storeVersion
	UseBestandsaufnahme   bool
	UseKonfigErworbenFile bool
	CompressionType       CompressionType
	LockInternalFiles     bool
}

func Default() Konfig {
	return Konfig{
		StoreVersion:          storeVersion(values.Int(3)),
		UseBestandsaufnahme:   true,
		UseKonfigErworbenFile: true,
		CompressionType:       CompressionTypeDefault,
		LockInternalFiles:     true,
	}
}

func (k Konfig) GetStoreVersion() schnittstellen.StoreVersion {
	return k.StoreVersion
}

func (k *Konfig) AddToFlagSet(f *flag.FlagSet) {
	f.BoolVar(
		&k.UseBestandsaufnahme,
		"use-bestandsaufnahme",
		k.UseBestandsaufnahme,
		"use bestandsaufnahme",
	)

	k.CompressionType.AddToFlagSet(f)

	f.BoolVar(
		&k.LockInternalFiles,
		"lock-internal-files",
		k.LockInternalFiles,
		"",
	)
}
