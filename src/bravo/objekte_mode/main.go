package objekte_mode

//go:generate stringer -type=Mode
type Mode int

const (
	ModeEmpty                 = Mode(iota)
	ModeAddToBestandsaufnahme = Mode(1 << iota)
	ModeUpdateTai
	ModeSchwanz

	ModeCommit = ModeAddToBestandsaufnahme | ModeUpdateTai | ModeSchwanz
)

func (a Mode) Contains(b Mode) bool {
	return a&b != 0
}
