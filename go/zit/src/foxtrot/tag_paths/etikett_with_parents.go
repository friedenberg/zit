package tag_paths

import (
	"fmt"
	"strings"
)

type EtikettWithParentsAndTypes struct {
	*Tag
	Parents PathsWithTypes
}

func (ewp EtikettWithParentsAndTypes) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s:%s", ewp.Tag, ewp.Parents)

	return sb.String()
}
