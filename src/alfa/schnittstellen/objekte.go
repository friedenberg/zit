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
	GetAkteSha() Sha
	GetObjekteSha() Sha
}

type StoredPtr interface {
	Stored
	SetAkteSha(Sha)
	SetObjekteSha(Sha)
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
