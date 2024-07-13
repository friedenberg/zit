package descriptions

type CliFormatTruncation int

const (
	CliFormatTruncationNone = CliFormatTruncation(iota)
	CliFormatTruncation66CharEllipsis
)
