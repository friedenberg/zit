package angeboren

import (
	"encoding/gob"
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
)

func init() {
	gob.RegisterName("KonfigLike", Konfig{})
}

type Konfig struct {
	StoreVersion                        storeVersion
	Recipients                          []string
	UseBestandsaufnahme                 bool
	UseKonfigErworbenFile               bool
	UseBestandsaufnahmeForVerzeichnisse bool
	CompressionType                     CompressionType
	LockInternalFiles                   bool
}

func Default() Konfig {
	return Konfig{
		StoreVersion:                        storeVersion(values.Int(5)),
		Recipients:                          make([]string, 0),
		UseBestandsaufnahme:                 true,
		UseBestandsaufnahmeForVerzeichnisse: true,
		UseKonfigErworbenFile:               true,
		CompressionType:                     CompressionTypeDefault,
		LockInternalFiles:                   true,
	}
}

func (k Konfig) GetStoreVersion() interfaces.StoreVersion {
	return k.StoreVersion
}

func (k Konfig) GetUseBestandsaufnahmeForVerzeichnisse() bool {
	return k.UseBestandsaufnahmeForVerzeichnisse
}

func (k *Konfig) AddToFlagSet(f *flag.FlagSet) {
	f.BoolVar(
		&k.UseBestandsaufnahme,
		"use-bestandsaufnahme",
		k.UseBestandsaufnahme,
		"use bestandsaufnahme",
	)

	f.BoolVar(
		&k.UseBestandsaufnahmeForVerzeichnisse,
		"use-bestandsaufnahme-for-verzeichnisse",
		k.UseBestandsaufnahmeForVerzeichnisse,
		"use bestandsaufnahme for verzeichnisse",
	)

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
