package query_spec

import (
	"strings"
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type tokenScannerTestCase struct {
	input    string
	expected []string
}

func getTokenScannerTestCases() []tokenScannerTestCase {
	return []tokenScannerTestCase{
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
		{
			input: "[uno/dos    !pdf     zz-inbox]",
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
		{
			input: `[
      /browser/bookmark-1FuOLQOYZAsP/ !toml-bookmark
      zz-site-org-mozilla-support
      !browser-bookmark "Get Help"
      url="https://support.\"mozilla.org/products/firefox"
      zz-site-org-mozilla-support] Get Help`,
			expected: []string{
				"[",
				"\n",
				"/browser/bookmark-1FuOLQOYZAsP/",
				" ",
				"!toml-bookmark",
				"\n",
				"zz-site-org-mozilla-support",
				"\n",
				"!browser-bookmark",
				" ",
				`"Get Help"`,
				"\n",
				`url="https://support.\"mozilla.org/products/firefox"`,
				"\n",
				"zz-site-org-mozilla-support",
				"]",
				" ",
				"Get",
				" ",
				"Help",
			},
		},
	}
}

func TestTokenScanner(t1 *testing.T) {
	t := test_logz.T{T: t1}

	var scanner TokenScanner

	for _, tc := range getTokenScannerTestCases() {
		scanner.Reset(strings.NewReader(tc.input))

		actual := make([]string, 0)

		for scanner.Scan() {
			t1 := scanner.GetToken().String()
			actual = append(actual, t1)
		}

		if err := scanner.Error(); err != nil {
			t.AssertNoError(err)
		}

		t.AssertNotEqual(tc.expected, actual)
	}
}

func TestTokenScannerWithTypes(t1 *testing.T) {
	t := test_logz.T{T: t1}

	input := `[
      /browser/bookmark-1FuOLQOYZAsP/ !toml-bookmark
      zz-site-org-mozilla-support
      !browser-bookmark "Get Help"
      url="https://support.\"mozilla.org/products/firefox"
      zz-site-org-mozilla-support] Get Help`

	type stringWithType struct {
		Value string
		TokenType
	}

	expected := []stringWithType{
		{Value: "[", TokenType: TokenTypeOperator},
		{Value: "\n", TokenType: TokenTypeOperator},
		{Value: "/browser/bookmark-1FuOLQOYZAsP/", TokenType: TokenTypeIdentifier},
		{Value: " ", TokenType: TokenTypeOperator},
		{Value: "!toml-bookmark", TokenType: TokenTypeIdentifier},
		{Value: "\n", TokenType: TokenTypeOperator},
		{Value: "zz-site-org-mozilla-support", TokenType: TokenTypeIdentifier},
		{Value: "\n", TokenType: TokenTypeOperator},
		{Value: "!browser-bookmark", TokenType: TokenTypeIdentifier},
		{Value: " ", TokenType: TokenTypeOperator},
		{Value: `"Get Help"`, TokenType: TokenTypeLiteral},
		{Value: "\n", TokenType: TokenTypeOperator},
		{Value: `url="https://support.\"mozilla.org/products/firefox"`, TokenType: TokenTypeField},
		{Value: "\n", TokenType: TokenTypeOperator},
		{Value: "zz-site-org-mozilla-support", TokenType: TokenTypeIdentifier},
		{Value: "]", TokenType: TokenTypeOperator},
		{Value: " ", TokenType: TokenTypeOperator},
		{Value: "Get", TokenType: TokenTypeIdentifier},
		{Value: " ", TokenType: TokenTypeOperator},
		{Value: "Help", TokenType: TokenTypeIdentifier},
	}

	var scanner TokenScanner
	scanner.Reset(strings.NewReader(input))

	actual := make([]stringWithType, 0)

	for scanner.Scan() {
		t1, ty := scanner.GetTokenAndType()
		actual = append(actual, stringWithType{t1.String(), ty})
	}

	if err := scanner.Error(); err != nil {
		t.AssertNoError(err)
	}

	t.AssertNotEqual(expected, actual)
}

func TestTokenScannerWithTypesAndParts(t1 *testing.T) {
	t := test_logz.T{T: t1}

	input := `[
      /browser/bookmark-1FuOLQOYZAsP/ !toml-bookmark
      zz-site-org-mozilla-support
      !browser-bookmark "Get Help"
      url="https://support.\"mozilla.org/products/firefox"
      zz-site-org-mozilla-support] Get Help`

	type stringWithTypeAndParts struct {
		Value string
		TokenType
		TokenParts
	}

	expected := []stringWithTypeAndParts{
		{Value: "[", TokenType: TokenTypeOperator},
		{Value: "\n", TokenType: TokenTypeOperator},
		{
			Value: "/browser/bookmark-1FuOLQOYZAsP/", TokenType: TokenTypeIdentifier, TokenParts: TokenParts{
				[]uint8("/browser/bookmark-1FuOLQOYZAsP/"),
				[]byte{},
			},
		},
		{Value: " ", TokenType: TokenTypeOperator},
		{Value: "!toml-bookmark", TokenType: TokenTypeIdentifier, TokenParts: TokenParts{
			[]uint8("!toml-bookmark"),
			[]byte{},
		}},
		{Value: "\n", TokenType: TokenTypeOperator},
		{Value: "zz-site-org-mozilla-support", TokenType: TokenTypeIdentifier, TokenParts: TokenParts{
			[]uint8("zz-site-org-mozilla-support"),
			[]byte{},
		}},
		{Value: "\n", TokenType: TokenTypeOperator},
		{Value: "!browser-bookmark", TokenType: TokenTypeIdentifier, TokenParts: TokenParts{
			[]uint8("!browser-bookmark"),
			[]byte{},
		}},
		{Value: " ", TokenType: TokenTypeOperator},
		{Value: `"Get Help"`, TokenType: TokenTypeLiteral, TokenParts: TokenParts{
			[]uint8("Get Help"),
			[]byte{},
		}},
		{Value: "\n", TokenType: TokenTypeOperator},
		{Value: `url="https://support.\"mozilla.org/products/firefox"`, TokenType: TokenTypeField, TokenParts: TokenParts{
			[]uint8("url"),
			[]uint8(`https://support.\"mozilla.org/products/firefox`),
		}},
		{Value: "\n", TokenType: TokenTypeOperator},
		{Value: "zz-site-org-mozilla-support", TokenType: TokenTypeIdentifier, TokenParts: TokenParts{
			[]uint8("zz-site-org-mozilla-support"),
			[]byte{},
		}},
		{Value: "]", TokenType: TokenTypeOperator},
		{Value: " ", TokenType: TokenTypeOperator},
		{Value: "Get", TokenType: TokenTypeIdentifier, TokenParts: TokenParts{
			[]uint8("Get"),
			[]byte{},
		}},
		{Value: " ", TokenType: TokenTypeOperator},
		{Value: "Help", TokenType: TokenTypeIdentifier, TokenParts: TokenParts{
			[]uint8("Help"),
			[]byte{},
		}},
	}

	var scanner TokenScanner
	scanner.Reset(strings.NewReader(input))

	actual := make([]stringWithTypeAndParts, 0)

	for scanner.Scan() {
		t1, ty, parts := scanner.GetTokenAndTypeAndParts()
		actual = append(actual, stringWithTypeAndParts{t1.String(), ty, parts.Clone()})
	}

	if err := scanner.Error(); err != nil {
		t.AssertNoError(err)
	}

	t.AssertNotEqual(expected, actual, cmpopts.EquateEmpty())
}

func TestTokenScannerWithTypesAndPartsRingBufferEdition(t1 *testing.T) {
	t := test_logz.T{T: t1}

	input := `[
      url="https://support.\"mozilla.org/products/firefox"
      zz-site-org-mozilla-support] Get Help`

	type stringWithTypeAndParts struct {
		Value string
		TokenType
		TokenParts
	}

	expected := []stringWithTypeAndParts{
		{Value: "[", TokenType: TokenTypeOperator},
		{Value: "\n", TokenType: TokenTypeOperator},
		{Value: `url="https://support.\"mozilla.org/products/firefox"`, TokenType: TokenTypeField, TokenParts: TokenParts{
			[]uint8("url"),
			[]uint8(`https://support.\"mozilla.org/products/firefox`),
		}},
		{Value: "\n", TokenType: TokenTypeOperator},
		{Value: "zz-site-org-mozilla-support", TokenType: TokenTypeIdentifier, TokenParts: TokenParts{
			[]uint8("zz-site-org-mozilla-support"),
			[]byte{},
		}},
		{Value: "]", TokenType: TokenTypeOperator},
		{Value: " ", TokenType: TokenTypeOperator},
		{Value: "Get", TokenType: TokenTypeIdentifier, TokenParts: TokenParts{
			[]uint8("Get"),
			[]byte{},
		}},
		{Value: " ", TokenType: TokenTypeOperator},
		{Value: "Help", TokenType: TokenTypeIdentifier, TokenParts: TokenParts{
			[]uint8("Help"),
			[]byte{},
		}},
	}

	var scanner TokenScanner
	rb := catgut.MakeRingBuffer(strings.NewReader(input), 0)
	rbrs := catgut.MakeRingBufferRuneScanner(rb)
	scanner.Reset(rbrs)
	// scanner.Reset(strings.NewReader(input))

	actual := make([]stringWithTypeAndParts, 0)

	for scanner.Scan() {
		t1, ty, parts := scanner.GetTokenAndTypeAndParts()
		actual = append(actual, stringWithTypeAndParts{t1.String(), ty, parts.Clone()})
	}

	if err := scanner.Error(); err != nil {
		t.AssertNoError(err)
	}

	t.AssertEqual(expected, actual, cmpopts.EquateEmpty())
}

func getTokenScannerTestCasesIdentifierLikeOnlySkipSpaces() []tokenScannerTestCase {
	return []tokenScannerTestCase{
		{
			input: "testing:e,t,k",
			expected: []string{
        "testing:e,t,k",
			},
		},
		{
			input: "[area-personal, area-work]:etikett",
			expected: []string{
				"[",
        "area-personal,",
        "area-work",
        "]",
        ":etikett",
			},
		},
		{
			input: " [ uno/dos ] bez",
			expected: []string{
				"[",
				"uno/dos",
				"]",
				"bez",
			},
		},
		{
			input: " [ uno/dos ] bez with spaces and more  spaces",
			expected: []string{
				"[",
				"uno/dos",
				"]",
				"bez",
				"with",
				"spaces",
				"and",
				"more",
				"spaces",
			},
		},
		{
			input: "[uno/dos    !pdf     zz-inbox]",
			expected: []string{
				"[",
				"uno/dos",
				"!pdf",
				"zz-inbox",
				"]",
			},
		},
		{
			input: `[
      /browser/bookmark-1FuOLQOYZAsP/ !toml-bookmark
      zz-site-org-mozilla-support
      !browser-bookmark "Get Help"
      url="https://support.\"mozilla.org/products/firefox"
      zz-site-org-mozilla-support] Get Help`,
			expected: []string{
				"[",
				"/browser/bookmark-1FuOLQOYZAsP/",
				"!toml-bookmark",
				"zz-site-org-mozilla-support",
				"!browser-bookmark",
				`"Get Help"`,
				`url="https://support.\"mozilla.org/products/firefox"`,
				"zz-site-org-mozilla-support",
				"]",
				"Get",
				"Help",
			},
		},
	}
}

func TestTokenScannerIdentifierLikeOnlySkipSpaces(t1 *testing.T) {
	t := test_logz.T{T: t1}

	var scanner TokenScanner

	for _, tc := range getTokenScannerTestCasesIdentifierLikeOnlySkipSpaces() {
		scanner.Reset(strings.NewReader(tc.input))

		actual := make([]string, 0)

		for scanner.ScanIdentifierLikeSkipSpaces() {
			t1 := scanner.GetToken().String()
			actual = append(actual, t1)
		}

		if err := scanner.Error(); err != nil {
			t.AssertNoError(err)
		}

		t.AssertNotEqual(tc.expected, actual)
	}
}
