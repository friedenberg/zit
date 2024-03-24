package query

import (
	"strings"
	"testing"

	"code.linenisgreat.com/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

func TestQuery(t1 *testing.T) {
	type testCase struct {
		stackInfo                                test_logz.StackInfo
		description, expected, expectedOptimized string
		defaultGattung                           kennung.Gattung
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
			expected: "!md:?Zettel",
			inputs:   []string{"!md?z"},
		},
		{
			expected: "ducks:Etikett [!md house]:+?Zettel",
			inputs:   []string{"!md?z", "house+z", "ducks:e"},
		},
		{
			expected: "ducks:Etikett [!md house]:?Zettel",
			inputs:   []string{"ducks:Etikett [!md house]?Zettel"},
		},
		{
			expected: "ducks:Etikett [=!md house]:?Zettel",
			inputs:   []string{"ducks:Etikett [=!md house]?Zettel"},
		},
		{
			expectedOptimized: "ducks:Etikett [=!md house wow]:?Zettel",
			expected:          "ducks:Etikett [=!md house wow]:?Zettel",
			inputs: []string{
				"ducks:Etikett [=!md house]?Zettel wow:Zettel",
			},
		},
		// { //TODO
		// 	stackInfo:         test_logz.MakeStackInfo(&t, 0),
		// 	expectedOptimized: "one/uno.zettel",
		// 	expected:          "one/uno.zettel",
		// 	inputs:            []string{"one/uno.zettel"},
		// },
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			expectedOptimized: "one/uno:Zettel",
			expected:          "one/uno:Zettel",
			defaultGattung:    kennung.MakeGattung(gattung.Zettel),
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
			defaultGattung:    kennung.MakeGattung(gattung.TrueGattung()...),
      expectedOptimized: ":Typ :Etikett :Zettel :Konfig :Kasten",
      expected:          ":Typ,Etikett,Zettel,Konfig,Kasten",
			inputs:            []string{},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			defaultGattung:    kennung.MakeGattung(gattung.TrueGattung()...),
			expectedOptimized: ":Typ :Etikett :Zettel :Konfig :Kasten",
			expected:          ":Typ,Etikett,Zettel,Konfig,Kasten",
			inputs:            []string{":"},
		},
	}

	for _, tc := range testCases {
		t1.Run(
			strings.Join(tc.inputs, " "),
			func(t1 *testing.T) {
				t := test_logz.T{T: t1}
				sut := (&Builder{}).WithDefaultGattungen(
					tc.defaultGattung,
				)

				m, err := sut.build(tc.inputs...)

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
