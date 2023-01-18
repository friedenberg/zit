package schnittstellen

type Gattung interface {
	GetGattungString() string
	GattungGetter
}

type GattungGetter interface {
	GetGattung() Gattung
}
