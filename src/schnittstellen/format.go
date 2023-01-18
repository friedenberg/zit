package schnittstellen

import "io"

//   _____                          _
//  |  ___|__  _ __ _ __ ___   __ _| |_
//  | |_ / _ \| '__| '_ ` _ \ / _` | __|
//  |  _| (_) | |  | | | | | | (_| | |_
//  |_|  \___/|_|  |_| |_| |_|\__,_|\__|
//

type FormatReader[T any, T1 Ptr[T]] interface {
	ReadFormat(io.Reader, T1) (int64, error)
}

type FormatWriter[T any, T1 Ptr[T]] interface {
	WriteFormat(io.Writer, T1) (int64, error)
}

type Parser[T any, T1 Ptr[T]] interface {
	Parse(io.Reader, T1) (int64, error)
}

type Formatter[T any, T1 Ptr[T]] interface {
	Format(io.Writer, T1) (int64, error)
}

type Format[T any, T1 Ptr[T]] interface {
	Parser[T, T1]
	Formatter[T, T1]
}
