package bezeichnung

type CliFormatTruncation int

const (
	CliFormatTruncationNone = CliFormatTruncation(iota)
	CliFormatTruncation66CharEllipsis
)
