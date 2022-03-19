package organize_text

type newZettel struct {
	bezeichnung string
}

// func (z newZettel) String() string {
// 	return fmt.Sprintf("- %s", z.bezeichnung)
// }

func (z *newZettel) Set(v string) (err error) {
	remaining := v

	if remaining[:2] != "- " {
		err = _Errorf("expected '- ', but got '%s'", remaining[:2])
		return
	}

	remaining = remaining[2:]

	z.bezeichnung = remaining

	return
}
