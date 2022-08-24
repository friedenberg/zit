package stored_zettel

import "fmt"

func (e External) String() string {
	return fmt.Sprintf("[%s %s]", e.Path, e.Stored.Sha)
}
