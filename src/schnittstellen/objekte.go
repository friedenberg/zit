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

type Stored interface {
	GattungGetter
	GetAkteSha() Sha
	GetObjekteSha() Sha
}

type StoredPtr interface {
	Stored
	SetAkteSha(Sha)
	SetObjekteSha(AkteReaderFactory, string) error
}

type VerzeichnissePtr[T any, T1 Objekte[T1]] interface {
	Ptr[T]
	ResetWithObjekte(T1)
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
