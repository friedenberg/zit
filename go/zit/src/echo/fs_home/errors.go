package fs_home

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type ErrNotInZitDir struct{}

func (e ErrNotInZitDir) Error() string {
	return "not in a zit directory"
}

func (e ErrNotInZitDir) ShouldShowStackTrace() bool {
	return false
}

func (e ErrNotInZitDir) Is(target error) (ok bool) {
	_, ok = target.(ErrNotInZitDir)
	return
}

func MakeErrAlreadyExists(
	sh interfaces.Sha,
	path string,
) (err *ErrAlreadyExists) {
	err = &ErrAlreadyExists{Path: path}
	err.Sha.SetShaLike(sh)
	return
}

type ErrAlreadyExists struct {
	sha.Sha
	Path string
}

func (e *ErrAlreadyExists) Error() string {
	return fmt.Sprintf("File with sha %s already exists: %s", &e.Sha, e.Path)
}

func (e *ErrAlreadyExists) Is(target error) bool {
	_, ok := target.(*ErrAlreadyExists)
	return ok
}
