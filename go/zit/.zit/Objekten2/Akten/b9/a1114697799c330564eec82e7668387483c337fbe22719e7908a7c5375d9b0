package errors

import (
	"fmt"
	"os"
)

type Signal struct {
	os.Signal
}

func (err Signal) Error() string {
	return fmt.Sprintf("signal: %q", err.Signal)
}
