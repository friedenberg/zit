package schnittstellen

type StringerGattungGetter interface {
	GattungGetter
	Stringer
}

type GattungLike interface {
	StringerGattungGetter
	EqualsGattung(GattungGetter) bool
	GetGattungString() string
	GetGattungStringPlural() string
}

type GattungGetter interface {
	GetGattung() GattungLike
}
