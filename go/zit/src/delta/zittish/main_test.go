package zittish

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"code.linenisgreat.com/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/src/charlie/catgut"
)

type testCase struct {
	input    string
	expected []string
}

func getTestCases() []testCase {
	return []testCase{
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
}

func TestZittish(t1 *testing.T) {
	t := test_logz.T{T: t1}

	for _, tc := range getTestCases() {
		scanner := catgut.NewScanner(
			catgut.MakeRingBuffer(strings.NewReader(tc.input), 0),
		)

		scanner.Split(SplitMatcher)

		actual := make([]*catgut.String, 0)

		for scanner.Scan() {
			t1 := catgut.Make(scanner.Text())
			actual = append(actual, t1)
		}

		if err := scanner.Err(); err != nil {
			t.AssertNoError(err)
		}

		expected := make([]*catgut.String, len(tc.expected))

		for i, v := range tc.expected {
			expected[i] = catgut.MakeFromString(v)
		}

		// t.Logf("%q", expected)

		if !reflect.DeepEqual(expected, actual) {
			t.NotEqual(tc.expected, actual)
		}
	}
}

func TestZittishRuneReader(t1 *testing.T) {
	t := test_logz.T{T: t1}

	for _, tc := range getTestCases() {
		rr := strings.NewReader(tc.input)

		actual := make([]string, 0)

		for {
			var token catgut.String
			err := NextToken(rr, &token)
			str := token.String()

			if err == nil && str == "" {
				t.Fatalf("expected non empty string")
			}

			if err == io.EOF {
				break
			} else {
				actual = append(actual, str)
				t.AssertNoError(err)
			}
		}

		if !reflect.DeepEqual(tc.expected, actual) {
			t.NotEqual(tc.expected, actual)
		}
	}
}
