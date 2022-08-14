package organize_text

import "github.com/friedenberg/zit/src/bravo/errors"

type newZettel struct {
	Bezeichnung string
}

// func (z newZettel) String() string {
// 	return fmt.Sprintf("- %s", z.bezeichnung)
// }

func (z *newZettel) Set(v string) (err error) {
	remaining := v

	if remaining[:2] != "- " {
		err = errors.Errorf("expected '- ', but got '%s'", remaining[:2])
		return
	}

	remaining = remaining[2:]

	z.Bezeichnung = remaining

	return
}
