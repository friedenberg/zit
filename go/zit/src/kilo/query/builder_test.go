package query

import (
	"strings"
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

func TestQuery(t1 *testing.T) {
	type testCase struct {
		stackInfo                                test_logz.StackInfo
		description, expected, expectedOptimized string
		defaultGattung                           ids.Genre
		inputs                                   []string
	}

	t := test_logz.T{T: t1}

	testCases := []testCase{
		{
			expected: "[[test,house] home]",
			inputs:   []string{"[test, house] home"},
		},
		{
			expected: "[[test,house] home wow]",
			inputs:   []string{"[test, house] home", "wow"},
		},
		{
			expected: "[^[test,house] home wow]",
			inputs:   []string{"^[test, house] home", "wow"},
		},
		{
			expected: "[[test,house] ^home wow]",
			inputs:   []string{"[test, house] ^home", "wow"},
		},
		{
			expected: "[[test,^house] home wow]",
			inputs:   []string{"[test, ^house] home", "wow"},
		},
		{
			expected: "[[test,house] home ^wow]",
			inputs:   []string{"[test, house] home", "^wow"},
		},
		{
			expected: "[^[[test,house] home] wow]",
			inputs:   []string{"^[[test, house] home]", "wow"},
		},
		{
			expected: "^[[test,house] home]:Zettel wow",
			inputs:   []string{"^[[test, house] home]:z", "wow"},
		},
		{
			expected: "[!md,home]:Zettel",
			inputs:   []string{"[!md,home]:z"},
		},
		{
			expected: "!md?Zettel",
			inputs:   []string{"!md?z"},
		},
		{
			expected: "ducks:Etikett [!md house]+?Zettel",
			inputs:   []string{"!md?z", "house+z", "ducks:e"},
		},
		{
			expected: "ducks:Etikett [!md house]?Zettel",
			inputs:   []string{"ducks:Etikett [!md house]?Zettel"},
		},
		{
			expected: "ducks:Etikett [=!md house]?Zettel",
			inputs:   []string{"ducks:Etikett [=!md house]?Zettel"},
		},
		{
			expectedOptimized: "ducks:Etikett [=!md house wow]:?Zettel",
			expected:          "ducks:Etikett [=!md house wow]:?Zettel",
			inputs: []string{
				"ducks:Etikett [=!md house]?Zettel wow:Zettel",
			},
		},
		{ // TODO try to make this expect `one/uno.zettel`
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			expectedOptimized: "one/uno:.Zettel",
			expected:          "one/uno:.Zettel",
			inputs:            []string{"one/uno.zettel"},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			expectedOptimized: "one/uno:Zettel",
			expected:          "one/uno:Zettel",
			defaultGattung:    ids.MakeGenre(genres.Zettel),
			inputs:            []string{"one/uno"},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			expectedOptimized: "one/uno:Zettel",
			expected:          "one/uno:Zettel",
			inputs:            []string{"one/uno:z"},
		},
		{
			expectedOptimized: ":Konfig",
			expected:          ":Konfig",
			inputs:            []string{":konfig"},
		},
		{
			expectedOptimized: ":Zettel",
			expected:          ":Zettel",
			inputs:            []string{":z"},
		},
		{
			expectedOptimized: ":Kasten",
			expected:          ":Kasten",
			inputs:            []string{":k"},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			expectedOptimized: "one/uno:+Zettel",
			expected:          "one/uno:+Zettel",
			inputs:            []string{"one/uno+"},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			expectedOptimized: "[one/dos, one/uno]:Zettel",
			expected:          "[one/dos, one/uno]:Zettel",
			inputs:            []string{"one/uno", "one/dos"},
		},
		{
			expectedOptimized: ":Typ :Etikett :Zettel",
			expected:          ":Typ,Etikett,Zettel",
			inputs:            []string{":z,t,e"},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			defaultGattung:    ids.MakeGenre(genres.TrueGenre()...),
			expectedOptimized: ":Akte :Typ :Etikett :Zettel :Konfig :Bestandsaufnahme :Kasten",
			expected:          ":Akte,Typ,Etikett,Zettel,Konfig,Bestandsaufnahme,Kasten",
			inputs:            []string{},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			defaultGattung:    ids.MakeGenre(genres.TrueGenre()...),
			expectedOptimized: ":Akte :Typ :Etikett :Zettel :Konfig :Bestandsaufnahme :Kasten",
			expected:          ":Akte,Typ,Etikett,Zettel,Konfig,Bestandsaufnahme,Kasten",
			inputs:            []string{":"},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			defaultGattung:    ids.MakeGenre(genres.TrueGenre()...),
			expectedOptimized: "2109504781.792086:Bestandsaufnahme",
			expected:          "2109504781.792086:Bestandsaufnahme",
			inputs:            []string{"[2109504781.792086]:b"},
		},
	}

	for _, tc := range testCases {
		t1.Run(
			strings.Join(tc.inputs, " "),
			func(t1 *testing.T) {
				t := test_logz.T{T: t1}
				sut := (&Builder{}).WithDefaultGenres(
					tc.defaultGattung,
				)

				m, err := sut.BuildQueryGroup(tc.inputs...)

				t.AssertNoError(err)
				actual := m.String()

				tcT := test_logz.TC{T: t, StackInfo: tc.stackInfo}

				if tc.expected != actual {
					tcT.Log("expected")
					tcT.AssertEqual(tc.expected, actual)
				}

				if tc.expectedOptimized == "" {
					return
				}

				actualOptimized := m.StringOptimized()

				if tc.expectedOptimized != actualOptimized {
					tcT.Log(m.StringDebug())
					tcT.Log("expectedOptimized")
					tcT.AssertEqual(tc.expectedOptimized, actualOptimized)
				}
			},
		)
	}
}
