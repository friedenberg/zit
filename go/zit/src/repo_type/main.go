package repo_type

//go:generate stringer -type=Type
type Type int

const (
	TypeUnknown = Type(iota)
	TypeReadWrite
	TypeOpaqueRelay
)
