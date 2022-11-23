package kennung

import "github.com/friedenberg/zit/src/alfa/errors"

func MatchTwoSortedEtikettStringSlices(a, b []string) (hasMatch bool) {
	var longer, shorter []string

	switch {
	case len(a) < len(b):
		shorter = a
		longer = b

	default:
		shorter = b
		longer = a
	}

	for _, v := range shorter {
		c := rune(v[0])

		var idx int
		idx, hasMatch = BinarySearchForRuneInEtikettenSortedStringSlice(longer, c)
		errors.Err().Print(idx, hasMatch)

		switch {
		case hasMatch:
			return

		case idx > len(longer)-1:
			return
		}

		longer = longer[idx:]
	}

	return
}

func BinarySearchForRuneInEtikettenSortedStringSlice(
	haystack []string,
	needle rune,
) (idx int, ok bool) {
	var low, hi int
	hi = len(haystack) - 1

	for {
		idx = ((hi - low) / 2) + low
		midValRaw := haystack[idx]

		if midValRaw == "" {
			return
		}

		midVal := rune(midValRaw[0])

		if hi == low {
			ok = midVal == needle
			return
		}

		switch {
		case midVal > needle:
			//search left
			hi = idx - 1
			continue

		case midVal == needle:
			//found
			ok = true
			return

		case midVal < needle:
			//search right
			low = idx + 1
			continue
		}
	}
}
