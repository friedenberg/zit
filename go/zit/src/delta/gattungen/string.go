package gattungen

import (
	"strings"

	"code.linenisgreat.com/zit/src/delta/gattung"
)

func String(gs Set) string {
	var sb strings.Builder
	first := true

	gs.Each(func(g gattung.Gattung) error {
		if !first {
			sb.WriteString(",")
		}

		sb.WriteString(g.String())
		first = false

		return nil
	})

	return sb.String()
}
