package etiketten_path

import (
	"fmt"
	"strings"
)

type EtikettWithParentsAndTypes struct {
	*Etikett
	Parents PathsWithTypes
}

func (ewp EtikettWithParentsAndTypes) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s:%s", ewp.Etikett, ewp.Parents)

	return sb.String()
}
