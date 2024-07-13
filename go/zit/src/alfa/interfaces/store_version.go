package interfaces

type StoreVersion interface {
	Stringer
	Lessor[StoreVersion]
	GetInt() int
}
