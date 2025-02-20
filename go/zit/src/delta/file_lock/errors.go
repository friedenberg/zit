package file_lock

import (
	"fmt"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
)

type ErrLockRequired struct {
	Operation string
}

func (e ErrLockRequired) Is(target error) bool {
	_, ok := target.(ErrLockRequired)
	return ok
}

func (e ErrLockRequired) Error() string {
	return fmt.Sprintf(
		"lock required for operation: %q",
		e.Operation,
	)
}

type ErrUnableToAcquireLock struct {
	envUI env_ui.Env
	Path  string
}

func (e ErrUnableToAcquireLock) Error() string {
	return fmt.Sprintf("repo is currently locked")
}

func (e ErrUnableToAcquireLock) Is(target error) bool {
	_, ok := target.(ErrUnableToAcquireLock)
	return ok
}

func (e ErrUnableToAcquireLock) GetHelpfulError() errors.Helpful {
	return e
}

func (e ErrUnableToAcquireLock) GetRetryableError() errors.Retryable {
	return e
}

func (e ErrUnableToAcquireLock) ErrorCause() []string {
	return []string{
		"A previous operation that acquired the repo lock failed.",
		"The lock is intentionally left behind in case recovery is necessary.",
	}
}

func (e ErrUnableToAcquireLock) ErrorRecovery() []string {
	return []string{
		fmt.Sprintf("The lockfile needs to removed (`rm %q`).", e.Path),
	}
}

func (err ErrUnableToAcquireLock) Recover(
	ctx errors.RetryableContext,
	in error,
) {
	errors.PrintHelpful(err.envUI.GetErr(), err)

	if err.envUI.Confirm("delete the existing lock?") {
		if err := os.Remove(err.Path); err != nil {
			ctx.CancelWithError(err)
		}

		ctx.Retry()
	} else {
		ctx.CancelWithBadRequestf("not deleting the lock. aborting.")
	}
}
