package matcher

import (
	"reflect"
	"testing"

	"github.com/friedenberg/zit/src/bravo/test_logz"
)

func TestBuilder(t1 *testing.T) {
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
	}

	for _, tc := range testCases {
		actual, err := getTokens(tc.input)
		if err != nil {
			t.Fatalf("expected no error but got %q", err)
		}

		if !reflect.DeepEqual(tc.expected, actual) {
			t.NotEqual(tc.expected, actual)
		}
	}
}
