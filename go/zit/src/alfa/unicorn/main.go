package unicorn

import "unicode/utf8"

func Not(f func(rune) bool) func(rune) bool {
	return func(r rune) bool {
		return !f(r)
	}
}

func CountRune(b []byte, r rune) (c int) {
	for i, w := 0, 0; i < len(b); i += w {
		runeValue, width := utf8.DecodeRune(b[i:])

		if runeValue != r {
			return
		}

		c++
		w = width
	}

	return
}
