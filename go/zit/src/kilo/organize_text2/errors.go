package organize_text2

type ErrorRead struct {
	error

	line, column int
}
