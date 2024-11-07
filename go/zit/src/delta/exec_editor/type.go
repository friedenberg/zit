package exec_editor

//go:generate stringer -type=Type
type Type byte

const (
	TypeUnknown = Type(iota)
	TypeVim
)
