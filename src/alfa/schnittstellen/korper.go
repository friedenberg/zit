package schnittstellen

type Korper interface {
	Element
	Stringer
	Kopf() string
	Schwanz() string
}
