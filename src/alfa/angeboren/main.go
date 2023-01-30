package angeboren

import (
	"encoding/gob"
	"flag"

	"github.com/friedenberg/zit/src/bravo/int_value"
)

type KonfigLike interface {
	GetStoreVersion() StoreVersion
}

type Getter interface {
	GetAngeboren() KonfigLike
}

func init() {
	gob.RegisterName("KonfigLike", Konfig{})
}

type Konfig struct {
	StoreVersion          storeVersion
	UseBestandsaufnahme   bool
	NewHinweisIndex       bool
	UseKonfigErworbenFile bool
}

func Default() Konfig {
	return Konfig{
		StoreVersion:          storeVersion(int_value.IntValue(1)),
		UseBestandsaufnahme:   true,
		UseKonfigErworbenFile: true,
		NewHinweisIndex:       true,
	}
}

func (k Konfig) GetStoreVersion() StoreVersion {
	return k.StoreVersion
}

func (k *Konfig) AddToFlags(f *flag.FlagSet) {
	f.BoolVar(&k.UseBestandsaufnahme, "use-bestandsaufnahme", k.UseBestandsaufnahme, "use bestandsaufnahme")
	f.BoolVar(&k.NewHinweisIndex, "use-new-hinweis-index", k.NewHinweisIndex, "use new hinweis index")
}
