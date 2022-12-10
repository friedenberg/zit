package age_io

import (
	"fmt"

	"github.com/friedenberg/zit/src/delta/sha"
)

type ErrAlreadyExists struct {
	sha.Sha
	Path string
}

func (e ErrAlreadyExists) Error() string {
	return fmt.Sprintf("File with sha %s already exists: %s", e.Sha, e.Path)
}

func (e ErrAlreadyExists) Is(target error) bool {
	_, ok := target.(ErrAlreadyExists)
	return ok
}
