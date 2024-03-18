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
		description, expected, expectedOptimized string
		defaultGattung                           kennung.Gattung
		inputs                                   []string
	}

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
			expectedOptimized: "ducks:Etikett [wow house !md]?Zettel",
			expected:          "ducks:Etikett [[=!md house] wow]?Zettel",
			inputs: []string{
				"ducks:Etikett [=!md house]?Zettel wow:Zettel",
			},
		},
		{
			expectedOptimized: "one/uno.zettel",
			expected:          "one/uno.zettel",
			inputs:            []string{"one/uno.zettel"},
		},
		{
			expectedOptimized: "one/uno:Zettel",
			expected:          "one/uno:Zettel",
			defaultGattung:    kennung.MakeGattung(gattung.Zettel),
			inputs:            []string{"one/uno"},
		},
		{
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
			expectedOptimized: "one/uno+",
			expected:          "one/uno+",
			inputs:            []string{"one/uno+"},
		},
		{
			expectedOptimized: "[one/uno, one/dos]",
			expected:          "[one/uno, one/dos]",
			inputs:            []string{"one/uno", "one/dos"},
		},
		{
			expectedOptimized: ":Typ :Etikett :Zettel",
			expected:          ":Typ,Etikett,Zettel",
			inputs:            []string{":z,t,e"},
		},
		{
			defaultGattung:    kennung.MakeGattung(gattung.TrueGattung()...),
			expectedOptimized: ":Typ :Etikett :Zettel",
			expected:          ":Typ,Etikett,Zettel",
			inputs:            []string{},
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

				if tc.expected != actual {
					// t.Logf("%#v", m)
					t.AssertEqual(tc.expected, actual)
				}

				if tc.expectedOptimized == "" {
					return
				}

				actualOptimized := m.StringOptimized()

				if tc.expectedOptimized != actualOptimized {
					t.Logf("%#v", m)
					t.Logf("%#v", m.OptimizedQueries[gattung.Zettel])
					t.Log(m.StringDebug())
					t.AssertEqual(tc.expectedOptimized, actualOptimized)
				}
			},
		)
	}
}
