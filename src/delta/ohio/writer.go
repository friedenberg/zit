package ohio

import "io"

func WriteAllOrDieTrying(w io.Writer, b []byte) (n int, err error) {
	var acc int

	for n < len(b) {
		acc, err = w.Write(b[n:])
		n += acc
		if err != nil {
			return
		}
	}

	return
}
