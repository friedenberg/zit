package zittish

import "unicode/utf8"

var mapMatcherOperators = map[rune]bool{
	' ': true,
	',': true,
	'{': true,
	'}': true,
	'[': true,
	']': true,
	':': true,
	'+': true,
	'.': true,
	'?': true,
}

func IsMatcherOperator(r rune) (ok bool) {
	_, ok = mapMatcherOperators[r]
	return
}

func SplitMatcher(
	data []byte,
	atEOF bool,
) (advance int, token []byte, err error) {
	for width, i := 0, 0; i < len(data); i += width {
		var r rune

		r, width = utf8.DecodeRune(data[i:])

		wasSplitRune := IsMatcherOperator(r)

		switch {
		case !wasSplitRune:
			continue

		case wasSplitRune && i == 0:
			return width, data[:width], nil

		default:
			return i, data[:i], nil
		}
	}

	// If we're at EOF, we have a final, non-empty, non-terminated word.  Return
	// it.
	if atEOF && len(data) > 0 {
		return len(data), data[0:], nil
	}

	return 0, nil, nil
}
