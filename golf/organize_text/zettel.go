package organize_text

import (
	"fmt"
	"strings"
)

type zettel struct {
	hinweis     string
	bezeichnung string
}

func (z zettel) String() string {
	return fmt.Sprintf("- [%s] %s", z.hinweis, z.bezeichnung)
}

func (z *zettel) Set(v string) (err error) {
	remaining := v

	if remaining[:3] != "- [" {
		err = _Errorf("expected '- [', but got '%s'", remaining[:2])
		return
	}

	remaining = remaining[3:]

	idx := -1

	if idx = strings.Index(remaining, "]"); idx == -1 {
		err = _Errorf("expected ']' after hinweis, but not found")
		return
	}

	z.hinweis = remaining[:idx]

	remaining = remaining[idx+1:]

	z.bezeichnung = remaining

	return
}
