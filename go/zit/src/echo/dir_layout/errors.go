package dir_layout

import (
	"errors"
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

type ErrBlobMissing struct {
	interfaces.ShaGetter
	Path string
}

func (e ErrBlobMissing) Error() string {
	return fmt.Sprintf(
		"Blob with sha %q does not exist locally: %q",
		e.GetShaLike(),
		e.Path,
	)
}

func (e ErrBlobMissing) Is(target error) bool {
	_, ok := target.(ErrBlobMissing)
	return ok
}

func IsErrBlobMissing(err error) bool {
	return errors.Is(err, ErrBlobMissing{})
}

func MakeErrTempAlreadyExists(
	path string,
) (err *ErrTempAlreadyExists) {
	err = &ErrTempAlreadyExists{Path: path}
	return
}

type ErrTempAlreadyExists struct {
	Path string
}

func (e *ErrTempAlreadyExists) Error() string {
	return fmt.Sprintf("Local temporary directory already exists: %q", e.Path)
}

func (e *ErrTempAlreadyExists) ErrorCause() string {
	return "Another zit previous process with the same PID likely terminated unexpectedly"
}

func (e *ErrTempAlreadyExists) ErrorRecovery() string {
	return "Check if there are any relevant files in the directory, or possible delete it"
}

func (e *ErrTempAlreadyExists) ErrorRecoveryAutomatic() string {
	return "TODO"
}

func (e *ErrTempAlreadyExists) Is(target error) bool {
	_, ok := target.(*ErrTempAlreadyExists)
	return ok
}
