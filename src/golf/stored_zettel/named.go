package stored_zettel

import "fmt"

func (z Named) String() string {
	return fmt.Sprintf("[%s %s]", z.Hinweis, z.Stored.Sha)
}
