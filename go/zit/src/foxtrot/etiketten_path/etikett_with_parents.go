package etiketten_path

import (
	"fmt"
	"strings"
)

type EtikettWithParents struct {
	*Etikett
	Parents SlicePaths
}

func (ewp EtikettWithParents) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s:%s", ewp.Etikett, ewp.Parents)

	return sb.String()
}
