package catgut

import "unicode/utf8"

func CompareUTF8Bytes(a, b []byte, partial bool) int {
	lenA, lenB := len(a), len(b)

	// TODO remove?
	switch {
	case lenA == 0 && lenB == 0:
		return 0

	case lenA == 0:
		return -1

	case lenB == 0:
		return 1
	}

	for {
		lenA, lenB := len(a), len(b)

		switch {
		case lenA == 0 && lenB == 0:
			return 0

		case lenA == 0:
			if partial && lenB <= lenA {
				return 0
			} else {
				return -1
			}

		case lenB == 0:
			if partial {
				return 0
			} else {
				return 1
			}
		}

		runeA, widthA := utf8.DecodeRune(a)
		a = a[widthA:]

		if runeA == utf8.RuneError {
			panic("not a valid utf8 string")
		}

		runeB, widthB := utf8.DecodeRune(b)
		b = b[widthB:]

		if runeB == utf8.RuneError {
			panic("not a valid utf8 string")
		}

		if runeA < runeB {
			return -1
		} else if runeA > runeB {
			return 1
		}
	}
}
