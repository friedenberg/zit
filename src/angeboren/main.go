package angeboren

import "flag"

type Konfig struct {
	UseBestandsaufnahme   bool
	UseKonfigErworbenFile bool
}

func Default() Konfig {
	return Konfig{
		UseBestandsaufnahme:   true,
		UseKonfigErworbenFile: true,
	}
}

func (k *Konfig) AddToFlags(f *flag.FlagSet) {
	f.BoolVar(&k.UseBestandsaufnahme, "use-bestandsaufnahme", false, "use bestandsaufnahme")
}
