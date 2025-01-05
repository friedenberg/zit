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
			stackInfo: test_logz.MakeStackInfo(&t, 0),
			expected:  "[[test,house] home]",
			inputs:    []string{"[test, house] home"},
		},
		{
			stackInfo: test_logz.MakeStackInfo(&t, 0),
			expected:  "[[test,house] home wow]",
			inputs:    []string{"[test, house] home", "wow"},
		},
		{
			stackInfo: test_logz.MakeStackInfo(&t, 0),
			expected:  "[^[test,house] home wow]",
			inputs:    []string{"^[test, house] home", "wow"},
		},
		{
			stackInfo: test_logz.MakeStackInfo(&t, 0),
			expected:  "[[test,house] ^home wow]",
			inputs:    []string{"[test, house] ^home", "wow"},
		},
		{
			stackInfo: test_logz.MakeStackInfo(&t, 0),
			expected:  "[[test,^house] home wow]",
			inputs:    []string{"[test, ^house] home", "wow"},
		},
		{
			stackInfo: test_logz.MakeStackInfo(&t, 0),
			expected:  "[[test,house] home ^wow]",
			inputs:    []string{"[test, house] home", "^wow"},
		},
		{
			stackInfo: test_logz.MakeStackInfo(&t, 0),
			expected:  "[^[[test,house] home] wow]",
			inputs:    []string{"^[[test, house] home]", "wow"},
		},
		{
			stackInfo: test_logz.MakeStackInfo(&t, 0),
			expected:  "^[[test,house] home]:Zettel wow",
			inputs:    []string{"^[[test, house] home]:z", "wow"},
		},
		{
			stackInfo: test_logz.MakeStackInfo(&t, 0),
			expected:  "[!md,home]:Zettel",
			inputs:    []string{"[!md,home]:z"},
		},
		{
			stackInfo: test_logz.MakeStackInfo(&t, 0),
			expected:  "!md?Zettel",
			inputs:    []string{"!md?z"},
		},
		{
			stackInfo: test_logz.MakeStackInfo(&t, 0),
			expected:  "ducks:Etikett [!md house]+?Zettel",
			inputs:    []string{"!md?z", "house+z", "ducks:e"},
		},
		{
			stackInfo: test_logz.MakeStackInfo(&t, 0),
			expected:  "ducks:Etikett [!md house]?Zettel",
			inputs:    []string{"ducks:Etikett [!md house]?Zettel"},
		},
		{
			stackInfo: test_logz.MakeStackInfo(&t, 0),
			expected:  "ducks:Etikett [=!md house]?Zettel",
			inputs:    []string{"ducks:Etikett [=!md house]?Zettel"},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
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
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			expectedOptimized: ":Konfig",
			expected:          ":Konfig",
			inputs:            []string{":konfig"},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			expectedOptimized: ":Zettel",
			expected:          ":Zettel",
			inputs:            []string{":z"},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
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
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			defaultGattung:    ids.MakeGenre(genres.TrueGenre()...),
			expectedOptimized: "^etikett-two.Zettel",
			expected:          "^etikett-two.Zettel",
			inputs:            []string{"^etikett-two.z"},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			defaultGattung:    ids.MakeGenre(genres.TrueGenre()...),
			expectedOptimized: "!md.Akte !md.Typ !md.Etikett !md.Zettel !md.Konfig !md.Bestandsaufnahme !md.Kasten",
			expected:          "!md.Akte,Typ,Etikett,Zettel,Konfig,Bestandsaufnahme,Kasten",
			inputs:            []string{"!md."},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			defaultGattung:    ids.MakeGenre(genres.TrueGenre()...),
			expectedOptimized: "-etikett-two.Zettel",
			expected:          "-etikett-two.Zettel",
			inputs:            []string{"-etikett-two.z"},
		},
	}

	for _, tc := range testCases {
		t1.Run(
			strings.Join(tc.inputs, " "),
			func(t1 *testing.T) {
				t := test_logz.TC{
					T:         test_logz.T{T: t1},
					StackInfo: tc.stackInfo,
				}

				sut := (&Builder{}).WithDefaultGenres(
					tc.defaultGattung,
				)

				m, err := sut.BuildQueryGroup(tc.inputs...)

				t.AssertNoError(err)
				actual := m.String()

				if tc.expected != actual {
					t.Log("expected")
					t.AssertEqual(tc.expected, actual)
				}

				if tc.expectedOptimized == "" {
					return
				}

				actualOptimized := m.StringOptimized()

				if tc.expectedOptimized != actualOptimized {
					t.Log(m.StringDebug())
					t.Log("expectedOptimized")
					t.AssertEqual(tc.expectedOptimized, actualOptimized)
				}
			},
		)
	}
}
