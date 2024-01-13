package schnittstellen

type StoreVersion interface {
	Stringer
	Lessor[StoreVersion]
	GetInt() int
}
