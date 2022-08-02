package hinweisen

import "github.com/friedenberg/zit/alfa/errors"

var (
	ErrDoesNotExist error
)

func init() {
	ErrDoesNotExist = errors.Normalf("hinweis does not exist")
}
