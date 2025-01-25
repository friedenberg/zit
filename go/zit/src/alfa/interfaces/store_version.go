package interfaces

type StoreVersion interface {
	Stringer
	GetInt() int
}
