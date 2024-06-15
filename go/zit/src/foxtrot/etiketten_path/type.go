package etiketten_path

type Type int

// describe these
const (
	TypeDirect = Type(iota)
	TypeInherit
	TypeIndirect
)
