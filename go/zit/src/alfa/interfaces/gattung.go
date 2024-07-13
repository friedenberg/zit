package interfaces

type StringerGattungGetter interface {
	GattungGetter
	Stringer
}

type StringerGattungKastenGetter interface {
	GattungGetter
	KastenGetter
	Stringer
}

type KastenLike interface {
	Stringer
	EqualsKasten(KastenGetter) bool
	GetKastenString() string
}

type KastenGetter interface {
	GetKasten() KastenLike
}

type GattungLike interface {
	StringerGattungGetter
	EqualsGattung(GattungGetter) bool
	GetGattungBitInt() byte
	GetGattungString() string
	GetGattungStringPlural() string
}

type GattungGetter interface {
	GetGattung() GattungLike
}
