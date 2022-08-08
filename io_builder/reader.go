package io_builder

import "io"

type ReaderFunc func(io.Reader) (io.Reader, error)

func Reader(rs ...ReaderFunc) (r io.Reader, err error) {
	return
}
