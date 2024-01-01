package objekte_update_type

type Type int

const (
	ModeEmpty                 = Type(iota)
	ModeAddToBestandsaufnahme = Type(1 << iota)
	ModeUpdateTai

	ModeCommit = ModeAddToBestandsaufnahme | ModeUpdateTai
)

func (a Type) Contains(b Type) bool {
	return a&b != 0
}
