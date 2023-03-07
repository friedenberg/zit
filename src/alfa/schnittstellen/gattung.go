package schnittstellen

type Gattung interface {
	Element
	GetGattungString() string
	GattungGetter
}

type GattungGetter interface {
	GetGattung() Gattung
}

type GattungenGetter interface {
	GetGattungen() Set[Gattung]
}
