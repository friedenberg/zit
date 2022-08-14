package hinweisen

import "github.com/friedenberg/zit/bravo/errors"

var (
	ErrDoesNotExist error
)

func init() {
	ErrDoesNotExist = errors.Normalf("hinweis does not exist")
}
