package schnittstellen

type Gattung interface {
	GetGattungString() string
	Equals(Gattung) bool
	GattungGetter
}

type GattungGetter interface {
	GetGattung() Gattung
}
