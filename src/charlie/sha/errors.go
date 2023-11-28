package sha

import "fmt"

type errLength [2]int

func makeErrLength(expected, actual int) error {
	if expected != actual {
		return errLength{expected, actual}
	} else {
		return nil
	}
}

func (e errLength) Error() string {
	return fmt.Sprintf("expected %d but got %d", e[0], e[1])
}
