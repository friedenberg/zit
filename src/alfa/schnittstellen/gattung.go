package schnittstellen

type Gattung interface {
	Element
	EqualsGattung(GattungGetter) bool
	GetGattungString() string
	GetGattungStringPlural() string
	GattungGetter
}

type GattungGetter interface {
	GetGattung() Gattung
}

type GattungenGetter interface {
	GetGattungen() Set[Gattung]
}
