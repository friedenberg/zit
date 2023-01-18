package schnittstellen

type Objekte[T any] interface {
	GattungGetter
	Equatable[T]
	GetAkteSha() Sha
}

type ObjektePtr[T any] interface {
	Objekte[T]
	Ptr[T]
	Resetable[T]
	SetAkteSha(Sha)
}
