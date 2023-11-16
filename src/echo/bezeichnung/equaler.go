package bezeichnung

var Equaler equaler

type equaler struct{}

func (equaler) Equals(a, b Bezeichnung) bool {
	return a.value == b.value
}

func (equaler) EqualsPtr(a, b *Bezeichnung) bool {
	return a.value == b.value
}
