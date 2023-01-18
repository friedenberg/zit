package schnittstellen

import "fmt"

type Korper interface {
	fmt.Stringer
	Kopf() string
	Schwanz() string
}
