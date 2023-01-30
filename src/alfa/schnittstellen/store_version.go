package schnittstellen

type StoreVersion interface {
	Stringer
	Lessor[StoreVersion]
	Int() int
}
