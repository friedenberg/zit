package schnittstellen

type Objekte[T any] interface {
	GattungGetter
	Equatable[T]
}

type ObjektePtr[T any] interface {
	Objekte[T]
	Ptr[T]
	Resetable[T]
}

type Stored interface {
	GattungGetter
	GetAkteSha() ShaLike
	GetObjekteSha() ShaLike
}

type StoredPtr interface {
	Stored
	SetAkteSha(ShaLike)
	SetObjekteSha(ShaLike)
}

type Transacted[T any] interface {
	Stored
	GetKennungString() string
}

type TransactedPtr[T any] interface {
	Transacted[T]
	Ptr[T]
	Resetable[T]
	StoredPtr
}
