package angeboren

import (
	"encoding/gob"
	"flag"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/int_value"
)

func init() {
	gob.RegisterName("KonfigLike", Konfig{})
}

type Konfig struct {
	StoreVersion          storeVersion
	UseBestandsaufnahme   bool
	UseKonfigErworbenFile bool
}

func Default() Konfig {
	return Konfig{
		StoreVersion:          storeVersion(int_value.IntValue(1)),
		UseBestandsaufnahme:   true,
		UseKonfigErworbenFile: true,
	}
}

func (k Konfig) GetStoreVersion() schnittstellen.StoreVersion {
	return k.StoreVersion
}

func (k *Konfig) AddToFlags(f *flag.FlagSet) {
	f.BoolVar(&k.UseBestandsaufnahme, "use-bestandsaufnahme", k.UseBestandsaufnahme, "use bestandsaufnahme")
}
