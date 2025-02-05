package env_workspace

import "code.linenisgreat.com/zit/go/zit/src/alfa/errors"

type ErrNotInWorkspace struct{
  *env
}

func (err ErrNotInWorkspace) Error() string {
	return "not in a workspace"
}

func (err ErrNotInWorkspace) Is(target error) bool {
	_, ok := target.(ErrNotInWorkspace)
	return ok
}

func (err ErrNotInWorkspace) ShouldShowStackTrace() bool {
	return false
}

// func (err ErrNotInWorkspace) ErrorCause() []string {
// }

// func (err ErrNotInWorkspace) ErrorRecovery() []string {
// }

func (err ErrNotInWorkspace) GetRetryableError() errors.Retryable {
	return err
}

func (err ErrNotInWorkspace) Recover(context errors.Context) {
}
