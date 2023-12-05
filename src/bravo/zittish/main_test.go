package zittish

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	"github.com/friedenberg/zit/src/bravo/test_logz"
)

func TestZittish(t1 *testing.T) {
	t := test_logz.T{T: t1}

	type testCase struct {
		input    string
		expected []string
	}

	testCases := []testCase{
		{
			input: "testing:e,t,k",
			expected: []string{
				"testing",
				":",
				"e",
				",",
				"t",
				",",
				"k",
			},
		},
		{
			input: "[area-personal, area-work]:etikett",
			expected: []string{
				"[",
				"area-personal",
				",",
				" ",
				"area-work",
				"]",
				":",
				"etikett",
			},
		},
		{
			input: " [ uno/dos ] bez",
			expected: []string{
				" ",
				"[",
				" ",
				"uno/dos",
				" ",
				"]",
				" ",
				"bez",
			},
		},
		{
			input: " [ uno/dos ] bez with spaces and more  spaces",
			expected: []string{
				" ",
				"[",
				" ",
				"uno/dos",
				" ",
				"]",
				" ",
				"bez",
				" ",
				"with",
				" ",
				"spaces",
				" ",
				"and",
				" ",
				"more",
				" ",
				" ",
				"spaces",
			},
		},
		{
			input: "[uno/dos !pdf zz-inbox]",
			expected: []string{
				"[",
				"uno/dos",
				" ",
				"!pdf",
				" ",
				"zz-inbox",
				"]",
			},
		},
	}

	for _, tc := range testCases {
		scanner := bufio.NewScanner(strings.NewReader(tc.input))

		scanner.Split(SplitMatcher)

		actual := make([]string, 0)

		for scanner.Scan() {
			actual = append(actual, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			t.AssertNoError(err)
		}

		if !reflect.DeepEqual(tc.expected, actual) {
			t.NotEqual(tc.expected, actual)
		}
	}
}
