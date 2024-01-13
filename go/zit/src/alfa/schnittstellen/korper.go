package schnittstellen

type Korper interface {
	Stringer
	Kopf() string
	Schwanz() string
}
