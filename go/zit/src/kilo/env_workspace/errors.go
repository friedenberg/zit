package env_workspace

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/mike/store_workspace_supplies"
)

type ErrUnsupportedType ids.Type

func (e ErrUnsupportedType) Is(target error) bool {
	_, ok := target.(ErrUnsupportedType)
	return ok
}

func (e ErrUnsupportedType) Error() string {
	return fmt.Sprintf("unsupported typ: %q", ids.Type(e))
}

func makeErrUnsupportedOperation(s *Store, op interface{}) error {
	return errors.WrapN(1,
		&ErrUnsupportedOperation{
			repoId:             s.RepoId,
			store:              s.StoreLike,
			operationInterface: op,
		},
	)
}

type ErrUnsupportedOperation struct {
	repoId             ids.RepoId
	store              store_workspace_supplies.StoreLike
	operationInterface interface{}
}

func (e ErrUnsupportedOperation) Error() string {
	return fmt.Sprintf(
		"store (%q:%T) does not support operation '%T'",
		e.repoId,
		e.store,
		e.operationInterface,
	)
}

func (e ErrUnsupportedOperation) Is(target error) bool {
	_, ok := target.(ErrUnsupportedOperation)
	return ok
}
