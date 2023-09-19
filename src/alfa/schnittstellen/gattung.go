package schnittstellen

type StringerGattungGetter interface {
	GattungGetter
	Stringer
}

type GattungLike interface {
	Element
	EqualsGattung(GattungGetter) bool
	GetGattungString() string
	GetGattungStringPlural() string
	GattungGetter
}

type GattungGetter interface {
	GetGattung() GattungLike
}

type GattungenGetter interface {
	GetGattungen() SetLike[GattungLike]
}
