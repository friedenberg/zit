package ts

type Clock interface {
	GetTime() Time
	GetTai() Tai
}
