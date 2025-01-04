package box_scanner

import (
	"fmt"
	"strings"
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/alfa/box"
	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
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
				`Get Help`,
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

type testToken struct {
	box.TokenType
	Contents string
}

func (token testToken) String() string {
	return fmt.Sprintf("%s %s", token.TokenType, token.Contents)
}

func makeTestToken(tt box.TokenType, contents string) testToken {
	return testToken{
		TokenType: tt,
		Contents:  contents,
	}
}

type testSeq []testToken

func makeTestSeq(tokens ...any) (ts testSeq) {
	for i := 0; i < len(tokens); i += 2 {
		ts = append(ts,
			makeTestToken(
				tokens[i].(box.TokenType),
				tokens[i+1].(string),
			),
		)
	}

	return
}

func makeTestSeqFromSeq(seq Seq) (ts testSeq) {
	for _, t := range seq {
		ts = append(ts, testToken{
			TokenType: t.TokenType,
			Contents:  string(t.Contents),
		})
	}

	return
}

type scannerTestCase struct {
	input    string
	expected []testSeq
}

func getScannerTestCases() []scannerTestCase {
	return []scannerTestCase{
		{
			input: "testing:e,t,k",
			expected: []testSeq{
				makeTestSeq(box.TokenTypeIdentifier, "testing"),
				makeTestSeq(box.TokenTypeOperator, ":"),
				makeTestSeq(box.TokenTypeIdentifier, "e"),
				makeTestSeq(box.TokenTypeOperator, ","),
				makeTestSeq(box.TokenTypeIdentifier, "t"),
				makeTestSeq(box.TokenTypeOperator, ","),
				makeTestSeq(box.TokenTypeIdentifier, "k"),
			},
		},
		{
			input: "[area-personal, area-work]:etikett",
			expected: []testSeq{
				makeTestSeq(box.TokenTypeOperator, "["),
				makeTestSeq(box.TokenTypeIdentifier, "area-personal"),
				makeTestSeq(box.TokenTypeOperator, ","),
				makeTestSeq(box.TokenTypeOperator, " "),
				makeTestSeq(box.TokenTypeIdentifier, "area-work"),
				makeTestSeq(box.TokenTypeOperator, "]"),
				makeTestSeq(box.TokenTypeOperator, ":"),
				makeTestSeq(box.TokenTypeIdentifier, "etikett"),
			},
		},
		{
			input: " [ uno/dos ] bez",
			expected: []testSeq{
				makeTestSeq(box.TokenTypeOperator, " "),
				makeTestSeq(box.TokenTypeOperator, "["),
				makeTestSeq(box.TokenTypeOperator, " "),
				makeTestSeq(
					box.TokenTypeIdentifier, "uno",
					box.TokenTypeOperator, "/",
					box.TokenTypeIdentifier, "dos",
				),
				makeTestSeq(box.TokenTypeOperator, " "),
				makeTestSeq(box.TokenTypeOperator, "]"),
				makeTestSeq(box.TokenTypeOperator, " "),
				makeTestSeq(box.TokenTypeIdentifier, "bez"),
			},
		},
		{
			input: "md.type",
			expected: []testSeq{
				makeTestSeq(
					box.TokenTypeIdentifier, "md",
					box.TokenTypeOperator, ".",
					box.TokenTypeIdentifier, "type",
				),
			},
		},
		{
			input: "[md.type]",
			expected: []testSeq{
				makeTestSeq(box.TokenTypeOperator, "["),
				makeTestSeq(
					box.TokenTypeIdentifier, "md",
					box.TokenTypeOperator, ".",
					box.TokenTypeIdentifier, "type",
				),
				makeTestSeq(box.TokenTypeOperator, "]"),
			},
		},
		{
			input: "[uno/dos !pdf zz-inbox]",
			expected: []testSeq{
				makeTestSeq(box.TokenTypeOperator, "["),
				makeTestSeq(
					box.TokenTypeIdentifier, "uno",
					box.TokenTypeOperator, "/",
					box.TokenTypeIdentifier, "dos",
				),
				makeTestSeq(box.TokenTypeOperator, " "),
				makeTestSeq(
					box.TokenTypeOperator, "!",
					box.TokenTypeIdentifier, "pdf",
				),
				makeTestSeq(box.TokenTypeOperator, " "),
				makeTestSeq(
					box.TokenTypeIdentifier, "zz-inbox",
				),
				makeTestSeq(box.TokenTypeOperator, "]"),
			},
		},
		{
			input: `/browser/bookmark-1FuOLQOYZAsP/ "Get Help" url="https://support.\"mozilla.org/products/firefox"`,
			expected: []testSeq{
				makeTestSeq(
					box.TokenTypeOperator, "/",
					box.TokenTypeIdentifier, "browser",
					box.TokenTypeOperator, "/",
					box.TokenTypeIdentifier, "bookmark-1FuOLQOYZAsP",
					box.TokenTypeOperator, "/",
				),
				makeTestSeq(box.TokenTypeOperator, " "),
				makeTestSeq(
					box.TokenTypeLiteral, "Get Help",
				),
				makeTestSeq(box.TokenTypeOperator, " "),
				makeTestSeq(
					box.TokenTypeIdentifier, "url",
					box.TokenTypeOperator, "=",
					box.TokenTypeLiteral, `https://support."mozilla.org/products/firefox`,
				),
			},
		},
	}
}

func TestTokenScanner(t1 *testing.T) {
	t := test_logz.T{T: t1}

	var scanner Scanner

	for _, tc := range getScannerTestCases() {
		scanner.Reset(strings.NewReader(tc.input))

		actual := make([]testSeq, 0)

		for scanner.Scan() {
			t1 := scanner.GetSeq().Clone()
			actual = append(actual, makeTestSeqFromSeq(t1))
		}

		if err := scanner.Error(); err != nil {
			t.AssertNoError(err)
		}

		t.Log(tc.input, "->", actual)

		t.AssertNotEqual(tc.expected, actual)
	}
}

// func TestTokenScannerWithTypes(t1 *testing.T) {
// 	t := test_logz.T{T: t1}

// 	input := `[
//       /browser/bookmark-1FuOLQOYZAsP/ !toml-bookmark
//       zz-site-org-mozilla-support
//       !browser-bookmark "Get Help"
//       url="https://support.\"mozilla.org/products/firefox"
//       zz-site-org-mozilla-support] Get Help`

// 	type stringWithType struct {
// 		Value     string
// 		TokenType box.TokenType
// 	}

// 	expected := []stringWithType{
// 		{Value: "[", TokenType: box.TokenTypeOperator},
// 		{Value: "\n", TokenType: box.TokenTypeOperator},
// 		{Value: "/browser/bookmark-1FuOLQOYZAsP/", TokenType: box.TokenTypeIdentifier},
// 		{Value: " ", TokenType: box.TokenTypeOperator},
// 		{Value: "!toml-bookmark", TokenType: box.TokenTypeIdentifier},
// 		{Value: "\n", TokenType: box.TokenTypeOperator},
// 		{Value: "zz-site-org-mozilla-support", TokenType: box.TokenTypeIdentifier},
// 		{Value: "\n", TokenType: box.TokenTypeOperator},
// 		{Value: "!browser-bookmark", TokenType: box.TokenTypeIdentifier},
// 		{Value: " ", TokenType: box.TokenTypeOperator},
// 		{Value: `"Get Help"`, TokenType: box.TokenTypeLiteral},
// 		{Value: "\n", TokenType: box.TokenTypeOperator},
// 		{Value: `url="https://support.\"mozilla.org/products/firefox"`, TokenType: box.TypeField},
// 		{Value: "\n", TokenType: box.TokenTypeOperator},
// 		{Value: "zz-site-org-mozilla-support", TokenType: box.TokenTypeIdentifier},
// 		{Value: "]", TokenType: box.TokenTypeOperator},
// 		{Value: " ", TokenType: box.TokenTypeOperator},
// 		{Value: "Get", TokenType: box.TokenTypeIdentifier},
// 		{Value: " ", TokenType: box.TokenTypeOperator},
// 		{Value: "Help", TokenType: box.TokenTypeIdentifier},
// 	}

// 	var scanner TokenScanner
// 	scanner.Reset(strings.NewReader(input))

// 	actual := make([]stringWithType, 0)

// 	for scanner.Scan() {
// 		t1, ty := scanner.GetTokenAndType()
// 		actual = append(actual, stringWithType{t1.String(), ty})
// 	}

// 	if err := scanner.Error(); err != nil {
// 		t.AssertNoError(err)
// 	}

// 	t.AssertNotEqual(expected, actual)
// }

// func TestTokenScannerWithTypesAndParts(t1 *testing.T) {
// 	t := test_logz.T{T: t1}

// 	input := `[
//       /browser/bookmark-1FuOLQOYZAsP/ !toml-bookmark
//       zz-site-org-mozilla-support
//       !browser-bookmark "Get Help"
//       url="https://support.\"mozilla.org/products/firefox"
//       zz-site-org-mozilla-support] Get Help`

// 	type stringWithTypeAndParts struct {
// 		Value     string
// 		TokenType box.TokenType
// 		Seq
// 	}

// 	expected := []stringWithTypeAndParts{
// 		{Value: "[", TokenType: box.TokenTypeOperator},
// 		{Value: "\n", TokenType: box.TokenTypeOperator},
// 		{
// 			Value: "/browser/bookmark-1FuOLQOYZAsP/", TokenType: box.TokenTypeIdentifier, Seq: Seq{
// 				[]uint8("/browser/bookmark-1FuOLQOYZAsP/"),
// 				[]byte{},
// 			},
// 		},
// 		{Value: " ", TokenType: box.TokenTypeOperator},
// 		{Value: "!toml-bookmark", TokenType: box.TokenTypeIdentifier, Seq: Seq{
// 			[]uint8("!toml-bookmark"),
// 			[]byte{},
// 		}},
// 		{Value: "\n", TokenType: box.TokenTypeOperator},
// 		{Value: "zz-site-org-mozilla-support", TokenType: box.TokenTypeIdentifier, Seq: Seq{
// 			[]uint8("zz-site-org-mozilla-support"),
// 			[]byte{},
// 		}},
// 		{Value: "\n", TokenType: box.TokenTypeOperator},
// 		{Value: "!browser-bookmark", TokenType: box.TokenTypeIdentifier, Seq: Seq{
// 			[]uint8("!browser-bookmark"),
// 			[]byte{},
// 		}},
// 		{Value: " ", TokenType: box.TokenTypeOperator},
// 		{Value: `"Get Help"`, TokenType: box.TokenTypeLiteral, Seq: Seq{
// 			[]uint8("Get Help"),
// 			[]byte{},
// 		}},
// 		{Value: "\n", TokenType: box.TokenTypeOperator},
// 		{Value: `url="https://support.\"mozilla.org/products/firefox"`, TokenType: box.TypeField, Seq: Seq{
// 			[]uint8("url"),
// 			[]uint8(`https://support.\"mozilla.org/products/firefox`),
// 		}},
// 		{Value: "\n", TokenType: box.TokenTypeOperator},
// 		{Value: "zz-site-org-mozilla-support", TokenType: box.TokenTypeIdentifier, Seq: Seq{
// 			[]uint8("zz-site-org-mozilla-support"),
// 			[]byte{},
// 		}},
// 		{Value: "]", TokenType: box.TokenTypeOperator},
// 		{Value: " ", TokenType: box.TokenTypeOperator},
// 		{Value: "Get", TokenType: box.TokenTypeIdentifier, Seq: Seq{
// 			[]uint8("Get"),
// 			[]byte{},
// 		}},
// 		{Value: " ", TokenType: box.TokenTypeOperator},
// 		{Value: "Help", TokenType: box.TokenTypeIdentifier, Seq: Seq{
// 			[]uint8("Help"),
// 			[]byte{},
// 		}},
// 	}

// 	var scanner TokenScanner
// 	scanner.Reset(strings.NewReader(input))

// 	actual := make([]stringWithTypeAndParts, 0)

// 	for scanner.Scan() {
// 		t1, ty, parts := scanner.GetTokenAndTypeAndParts()
// 		actual = append(actual, stringWithTypeAndParts{t1.String(), ty, parts.Clone()})
// 	}

// 	if err := scanner.Error(); err != nil {
// 		t.AssertNoError(err)
// 	}

// 	t.AssertNotEqual(expected, actual, cmpopts.EquateEmpty())
// }

// func TestTokenScannerWithTypesAndPartsRingBufferEdition(t1 *testing.T) {
// 	t := test_logz.T{T: t1}

// 	input := `[
//       url="https://support.\"mozilla.org/products/firefox"
//       zz-site-org-mozilla-support] Get Help`

// 	type stringWithTypeAndParts struct {
// 		Value string
// 		box.TokenType
// 		Seq
// 	}

// 	expected := []stringWithTypeAndParts{
// 		{Value: "[", TokenType: box.TokenTypeOperator},
// 		{Value: "\n", TokenType: box.TokenTypeOperator},
// 		{Value: `url="https://support.\"mozilla.org/products/firefox"`, TokenType: box.TypeField, Seq: Seq{
// 			[]uint8("url"),
// 			[]uint8(`https://support.\"mozilla.org/products/firefox`),
// 		}},
// 		{Value: "\n", TokenType: box.TokenTypeOperator},
// 		{Value: "zz-site-org-mozilla-support", TokenType: box.TokenTypeIdentifier, Seq: Seq{
// 			[]uint8("zz-site-org-mozilla-support"),
// 			[]byte{},
// 		}},
// 		{Value: "]", TokenType: box.TokenTypeOperator},
// 		{Value: " ", TokenType: box.TokenTypeOperator},
// 		{Value: "Get", TokenType: box.TokenTypeIdentifier, Seq: Seq{
// 			[]uint8("Get"),
// 			[]byte{},
// 		}},
// 		{Value: " ", TokenType: box.TokenTypeOperator},
// 		{Value: "Help", TokenType: box.TokenTypeIdentifier, Seq: Seq{
// 			[]uint8("Help"),
// 			[]byte{},
// 		}},
// 	}

// 	var scanner TokenScanner
// 	rb := catgut.MakeRingBuffer(strings.NewReader(input), 0)
// 	rbrs := catgut.MakeRingBufferRuneScanner(rb)
// 	scanner.Reset(rbrs)
// 	// scanner.Reset(strings.NewReader(input))

// 	actual := make([]stringWithTypeAndParts, 0)

// 	for scanner.Scan() {
// 		t1, ty, parts := scanner.GetTokenAndTypeAndParts()
// 		actual = append(actual, stringWithTypeAndParts{t1.String(), ty, parts.Clone()})
// 	}

// 	if err := scanner.Error(); err != nil {
// 		t.AssertNoError(err)
// 	}

// 	t.AssertEqual(expected, actual, cmpopts.EquateEmpty())
// }

type typeAndParts struct {
	TokenType          box.TokenType
	Token, Left, Right string
}

type tokenScannerTypesAndPartsTestCase struct {
	input    string
	expected []typeAndParts
}

// func getTokenScannerTypeAndPartsTestCases() []tokenScannerTypesAndPartsTestCase {
// 	return []tokenScannerTypesAndPartsTestCase{
// 		{
// 			input: `[/firefox-ddog/bookmark-5nSmpin9cwMc title="Equipment Recommendations" url="https://atlassian.net/"] Equipment Recommendations`,
// 			expected: []typeAndParts{
// 				{
// 					TokenType: box.TokenTypeOperator,
// 					Token:     "[",
// 				},
// 				{
// 					TokenType: box.TokenTypeIdentifier,
// 					Token:     "/firefox-ddog/bookmark-5nSmpin9cwMc",
// 					Left:      "/firefox-ddog/bookmark-5nSmpin9cwMc",
// 				},
// 				{
// 					TokenType: box.TokenTypeOperator,
// 					Token:     " ",
// 				},
// 				{
// 					TokenType: box.TypeField,
// 					Token:     `title="Equipment Recommendations"`,
// 					Left:      `title`,
// 					Right:     `Equipment Recommendations`,
// 				},
// 				{
// 					TokenType: box.TokenTypeOperator,
// 					Token:     " ",
// 				},
// 				{
// 					TokenType: box.TypeField,
// 					Token:     `url="https://atlassian.net/"`,
// 					Left:      `url`,
// 					Right:     `https://atlassian.net/`,
// 				},
// 				{
// 					TokenType: box.TokenTypeOperator,
// 					Token:     "]",
// 				},
// 				{
// 					TokenType: box.TokenTypeOperator,
// 					Token:     " ",
// 				},
// 				{
// 					TokenType: box.TokenTypeIdentifier,
// 					Token:     "Equipment",
// 					Left:      "Equipment",
// 				},
// 				{
// 					TokenType: box.TokenTypeOperator,
// 					Token:     " ",
// 				},
// 				{
// 					TokenType: box.TokenTypeIdentifier,
// 					Token:     "Recommendations",
// 					Left:      "Recommendations",
// 				},
// 			},
// 		},
// 	}
// }

// func TestTokenScannerTypesAndParts(t1 *testing.T) {
// 	t := test_logz.T{T: t1}

// 	var scanner TokenScanner

// 	for _, tc := range getTokenScannerTypeAndPartsTestCases() {
// 		scanner.Reset(strings.NewReader(tc.input))

// 		actual := make([]typeAndParts, 0)

// 		for scanner.Scan() {
// 			token, tokenType, parts := scanner.GetTokenAndTypeAndParts()
// 			actual = append(
// 				actual,
// 				typeAndParts{
// 					TokenType: tokenType,
// 					Token:     token.String(),
// 					Left:      string(parts.Left),
// 					Right:     string(parts.Right),
// 				},
// 			)
// 		}

// 		if err := scanner.Error(); err != nil {
// 			t.AssertNoError(err)
// 		}

// 		t.AssertNotEqual(tc.expected, actual)
// 	}
// }

// func getTokenScannerTypeAndPartsTestCasesSkipWhitespace() []tokenScannerTypesAndPartsTestCase {
// 	return []tokenScannerTypesAndPartsTestCase{
// 		{
// 			input: `[/firefox-ddog/bookmark-5nSmpin9cwMc title="Equipment Recommendations" url="https://atlassian.net/"] Equipment Recommendations`,
// 			expected: []typeAndParts{
// 				{
// 					TokenType: box.TokenTypeOperator,
// 					Token:     "[",
// 				},
// 				{
// 					TokenType: box.TokenTypeIdentifier,
// 					Token:     "/firefox-ddog/bookmark-5nSmpin9cwMc",
// 					Left:      "/firefox-ddog/bookmark-5nSmpin9cwMc",
// 				},
// 				{
// 					TokenType: box.TypeField,
// 					Token:     `title="Equipment Recommendations"`,
// 					Left:      `title`,
// 					Right:     `Equipment Recommendations`,
// 				},
// 				{
// 					TokenType: box.TypeField,
// 					Token:     `url="https://atlassian.net/"`,
// 					Left:      `url`,
// 					Right:     `https://atlassian.net/`,
// 				},
// 				{
// 					TokenType: box.TokenTypeOperator,
// 					Token:     "]",
// 				},
// 				{
// 					TokenType: box.TokenTypeIdentifier,
// 					Token:     "Equipment",
// 					Left:      "Equipment",
// 				},
// 				{
// 					TokenType: box.TokenTypeIdentifier,
// 					Token:     "Recommendations",
// 					Left:      "Recommendations",
// 				},
// 			},
// 		},
// 	}
// }

// func TestTokenScannerTypesAndPartsSkipWhitespace(t1 *testing.T) {
// 	t := test_logz.T{T: t1}

// 	var scanner TokenScanner

// 	for _, tc := range getTokenScannerTypeAndPartsTestCasesSkipWhitespace() {
// 		scanner.Reset(strings.NewReader(tc.input))

// 		actual := make([]typeAndParts, 0)

// 		for scanner.ScanSkipSpace() {
// 			token, tokenType, parts := scanner.GetTokenAndTypeAndParts()
// 			actual = append(
// 				actual,
// 				typeAndParts{
// 					TokenType: tokenType,
// 					Token:     token.String(),
// 					Left:      string(parts.Left),
// 					Right:     string(parts.Right),
// 				},
// 			)
// 		}

// 		if err := scanner.Error(); err != nil {
// 			t.AssertNoError(err)
// 		}

// 		t.AssertNotEqual(tc.expected, actual)
// 	}
// }
