package query_spec

import (
	"strings"
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/alfa/token_types"
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
		token_types.TokenType
	}

	expected := []stringWithType{
		{Value: "[", TokenType: token_types.TypeOperator},
		{Value: "\n", TokenType: token_types.TypeOperator},
		{Value: "/browser/bookmark-1FuOLQOYZAsP/", TokenType: token_types.TypeIdentifier},
		{Value: " ", TokenType: token_types.TypeOperator},
		{Value: "!toml-bookmark", TokenType: token_types.TypeIdentifier},
		{Value: "\n", TokenType: token_types.TypeOperator},
		{Value: "zz-site-org-mozilla-support", TokenType: token_types.TypeIdentifier},
		{Value: "\n", TokenType: token_types.TypeOperator},
		{Value: "!browser-bookmark", TokenType: token_types.TypeIdentifier},
		{Value: " ", TokenType: token_types.TypeOperator},
		{Value: `"Get Help"`, TokenType: token_types.TypeLiteral},
		{Value: "\n", TokenType: token_types.TypeOperator},
		{Value: `url="https://support.\"mozilla.org/products/firefox"`, TokenType: token_types.TypeField},
		{Value: "\n", TokenType: token_types.TypeOperator},
		{Value: "zz-site-org-mozilla-support", TokenType: token_types.TypeIdentifier},
		{Value: "]", TokenType: token_types.TypeOperator},
		{Value: " ", TokenType: token_types.TypeOperator},
		{Value: "Get", TokenType: token_types.TypeIdentifier},
		{Value: " ", TokenType: token_types.TypeOperator},
		{Value: "Help", TokenType: token_types.TypeIdentifier},
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
		token_types.TokenType
		TokenParts
	}

	expected := []stringWithTypeAndParts{
		{Value: "[", TokenType: token_types.TypeOperator},
		{Value: "\n", TokenType: token_types.TypeOperator},
		{
			Value: "/browser/bookmark-1FuOLQOYZAsP/", TokenType: token_types.TypeIdentifier, TokenParts: TokenParts{
				[]uint8("/browser/bookmark-1FuOLQOYZAsP/"),
				[]byte{},
			},
		},
		{Value: " ", TokenType: token_types.TypeOperator},
		{Value: "!toml-bookmark", TokenType: token_types.TypeIdentifier, TokenParts: TokenParts{
			[]uint8("!toml-bookmark"),
			[]byte{},
		}},
		{Value: "\n", TokenType: token_types.TypeOperator},
		{Value: "zz-site-org-mozilla-support", TokenType: token_types.TypeIdentifier, TokenParts: TokenParts{
			[]uint8("zz-site-org-mozilla-support"),
			[]byte{},
		}},
		{Value: "\n", TokenType: token_types.TypeOperator},
		{Value: "!browser-bookmark", TokenType: token_types.TypeIdentifier, TokenParts: TokenParts{
			[]uint8("!browser-bookmark"),
			[]byte{},
		}},
		{Value: " ", TokenType: token_types.TypeOperator},
		{Value: `"Get Help"`, TokenType: token_types.TypeLiteral, TokenParts: TokenParts{
			[]uint8("Get Help"),
			[]byte{},
		}},
		{Value: "\n", TokenType: token_types.TypeOperator},
		{Value: `url="https://support.\"mozilla.org/products/firefox"`, TokenType: token_types.TypeField, TokenParts: TokenParts{
			[]uint8("url"),
			[]uint8(`https://support.\"mozilla.org/products/firefox`),
		}},
		{Value: "\n", TokenType: token_types.TypeOperator},
		{Value: "zz-site-org-mozilla-support", TokenType: token_types.TypeIdentifier, TokenParts: TokenParts{
			[]uint8("zz-site-org-mozilla-support"),
			[]byte{},
		}},
		{Value: "]", TokenType: token_types.TypeOperator},
		{Value: " ", TokenType: token_types.TypeOperator},
		{Value: "Get", TokenType: token_types.TypeIdentifier, TokenParts: TokenParts{
			[]uint8("Get"),
			[]byte{},
		}},
		{Value: " ", TokenType: token_types.TypeOperator},
		{Value: "Help", TokenType: token_types.TypeIdentifier, TokenParts: TokenParts{
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
		token_types.TokenType
		TokenParts
	}

	expected := []stringWithTypeAndParts{
		{Value: "[", TokenType: token_types.TypeOperator},
		{Value: "\n", TokenType: token_types.TypeOperator},
		{Value: `url="https://support.\"mozilla.org/products/firefox"`, TokenType: token_types.TypeField, TokenParts: TokenParts{
			[]uint8("url"),
			[]uint8(`https://support.\"mozilla.org/products/firefox`),
		}},
		{Value: "\n", TokenType: token_types.TypeOperator},
		{Value: "zz-site-org-mozilla-support", TokenType: token_types.TypeIdentifier, TokenParts: TokenParts{
			[]uint8("zz-site-org-mozilla-support"),
			[]byte{},
		}},
		{Value: "]", TokenType: token_types.TypeOperator},
		{Value: " ", TokenType: token_types.TypeOperator},
		{Value: "Get", TokenType: token_types.TypeIdentifier, TokenParts: TokenParts{
			[]uint8("Get"),
			[]byte{},
		}},
		{Value: " ", TokenType: token_types.TypeOperator},
		{Value: "Help", TokenType: token_types.TypeIdentifier, TokenParts: TokenParts{
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

type typeAndParts struct {
	TokenType          token_types.TokenType
	Token, Left, Right string
}

type tokenScannerTypesAndPartsTestCase struct {
	input    string
	expected []typeAndParts
}

func getTokenScannerTypeAndPartsTestCases() []tokenScannerTypesAndPartsTestCase {
	return []tokenScannerTypesAndPartsTestCase{
		{
			input: `[/firefox-ddog/bookmark-5nSmpin9cwMc title="Equipment Recommendations" url="https://atlassian.net/"] Equipment Recommendations`,
			expected: []typeAndParts{
				{
					TokenType: token_types.TypeOperator,
					Token:     "[",
				},
				{
					TokenType: token_types.TypeIdentifier,
					Token:     "/firefox-ddog/bookmark-5nSmpin9cwMc",
					Left:      "/firefox-ddog/bookmark-5nSmpin9cwMc",
				},
				{
					TokenType: token_types.TypeOperator,
					Token:     " ",
				},
				{
					TokenType: token_types.TypeField,
					Token:     `title="Equipment Recommendations"`,
					Left:      `title`,
					Right:     `Equipment Recommendations`,
				},
				{
					TokenType: token_types.TypeOperator,
					Token:     " ",
				},
				{
					TokenType: token_types.TypeField,
					Token:     `url="https://atlassian.net/"`,
					Left:      `url`,
					Right:     `https://atlassian.net/`,
				},
				{
					TokenType: token_types.TypeOperator,
					Token:     "]",
				},
				{
					TokenType: token_types.TypeOperator,
					Token:     " ",
				},
				{
					TokenType: token_types.TypeIdentifier,
					Token:     "Equipment",
					Left:      "Equipment",
				},
				{
					TokenType: token_types.TypeOperator,
					Token:     " ",
				},
				{
					TokenType: token_types.TypeIdentifier,
					Token:     "Recommendations",
					Left:      "Recommendations",
				},
			},
		},
	}
}

func TestTokenScannerTypesAndParts(t1 *testing.T) {
	t := test_logz.T{T: t1}

	var scanner TokenScanner

	for _, tc := range getTokenScannerTypeAndPartsTestCases() {
		scanner.Reset(strings.NewReader(tc.input))

		actual := make([]typeAndParts, 0)

		for scanner.Scan() {
			token, tokenType, parts := scanner.GetTokenAndTypeAndParts()
			actual = append(
				actual,
				typeAndParts{
					TokenType: tokenType,
					Token:     token.String(),
					Left:      string(parts.Left),
					Right:     string(parts.Right),
				},
			)
		}

		if err := scanner.Error(); err != nil {
			t.AssertNoError(err)
		}

		t.AssertNotEqual(tc.expected, actual)
	}
}

func getTokenScannerTypeAndPartsTestCasesSkipWhitespace() []tokenScannerTypesAndPartsTestCase {
	return []tokenScannerTypesAndPartsTestCase{
		{
			input: `[/firefox-ddog/bookmark-5nSmpin9cwMc title="Equipment Recommendations" url="https://atlassian.net/"] Equipment Recommendations`,
			expected: []typeAndParts{
				{
					TokenType: token_types.TypeOperator,
					Token:     "[",
				},
				{
					TokenType: token_types.TypeIdentifier,
					Token:     "/firefox-ddog/bookmark-5nSmpin9cwMc",
					Left:      "/firefox-ddog/bookmark-5nSmpin9cwMc",
				},
				{
					TokenType: token_types.TypeField,
					Token:     `title="Equipment Recommendations"`,
					Left:      `title`,
					Right:     `Equipment Recommendations`,
				},
				{
					TokenType: token_types.TypeField,
					Token:     `url="https://atlassian.net/"`,
					Left:      `url`,
					Right:     `https://atlassian.net/`,
				},
				{
					TokenType: token_types.TypeOperator,
					Token:     "]",
				},
				{
					TokenType: token_types.TypeIdentifier,
					Token:     "Equipment",
					Left:      "Equipment",
				},
				{
					TokenType: token_types.TypeIdentifier,
					Token:     "Recommendations",
					Left:      "Recommendations",
				},
			},
		},
	}
}

func TestTokenScannerTypesAndPartsSkipWhitespace(t1 *testing.T) {
	t := test_logz.T{T: t1}

	var scanner TokenScanner

	for _, tc := range getTokenScannerTypeAndPartsTestCasesSkipWhitespace() {
		scanner.Reset(strings.NewReader(tc.input))

		actual := make([]typeAndParts, 0)

		for scanner.ScanSkipSpace() {
			token, tokenType, parts := scanner.GetTokenAndTypeAndParts()
			actual = append(
				actual,
				typeAndParts{
					TokenType: tokenType,
					Token:     token.String(),
					Left:      string(parts.Left),
					Right:     string(parts.Right),
				},
			)
		}

		if err := scanner.Error(); err != nil {
			t.AssertNoError(err)
		}

		t.AssertNotEqual(tc.expected, actual)
	}
}
