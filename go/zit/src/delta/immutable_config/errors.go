package immutable_config

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type ErrFutureStoreVersion struct {
	interfaces.StoreVersion
}

func (e ErrFutureStoreVersion) Error() string {
	return fmt.Sprintf(
		strings.Join(
			[]string{
				"store version is from the future: %q",
				"This means that this installation of zit is likely out of date.",
			},
			". ",
		),
		e.StoreVersion,
	)
}

func (e ErrFutureStoreVersion) Is(target error) bool {
	_, ok := target.(ErrFutureStoreVersion)
	return ok
}
