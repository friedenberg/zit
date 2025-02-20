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
		defaultGenre                             ids.Genre
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
			expected:  "ducks:Tag [!md house]+?Zettel",
			inputs:    []string{"!md?z", "house+z", "ducks:e"},
		},
		{
			stackInfo: test_logz.MakeStackInfo(&t, 0),
			expected:  "ducks:Tag [!md house]?Zettel",
			inputs:    []string{"ducks:Tag [!md house]?Zettel"},
		},
		{
			stackInfo: test_logz.MakeStackInfo(&t, 0),
			expected:  "ducks:Tag [=!md house]?Zettel",
			inputs:    []string{"ducks:Tag [=!md house]?Zettel"},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			expectedOptimized: "ducks:Tag [=!md house wow]:?Zettel",
			expected:          "ducks:Tag [=!md house wow]:?Zettel",
			inputs: []string{
				"ducks:Tag [=!md house]?Zettel wow:Zettel",
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
			defaultGenre:      ids.MakeGenre(genres.Zettel),
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
			expectedOptimized: ":Config",
			expected:          ":Config",
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
			expectedOptimized: ":Repo",
			expected:          ":Repo",
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
			expectedOptimized: ":Type :Tag :Zettel",
			expected:          ":Type,Tag,Zettel",
			inputs:            []string{":z,t,e"},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			defaultGenre:      ids.MakeGenre(genres.All()...),
			expectedOptimized: ":Blob :Type :Tag :Zettel :Config :InventoryList :Repo",
			expected:          ":Blob,Type,Tag,Zettel,Config,InventoryList,Repo",
			inputs:            []string{},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			defaultGenre:      ids.MakeGenre(genres.All()...),
			expectedOptimized: ":Blob :Type :Tag :Zettel :Config :InventoryList :Repo",
			expected:          ":Blob,Type,Tag,Zettel,Config,InventoryList,Repo",
			inputs:            []string{":"},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			defaultGenre:      ids.MakeGenre(genres.All()...),
			expectedOptimized: "2109504781.792086:InventoryList",
			expected:          "2109504781.792086:InventoryList",
			inputs:            []string{"[2109504781.792086]:b"},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			defaultGenre:      ids.MakeGenre(genres.All()...),
			expectedOptimized: "^etikett-two.Zettel",
			expected:          "^etikett-two.Zettel",
			inputs:            []string{"^etikett-two.z"},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			defaultGenre:      ids.MakeGenre(genres.All()...),
			expectedOptimized: "!md.Blob !md.Type !md.Tag !md.Zettel !md.Config !md.InventoryList !md.Repo",
			expected:          "!md.Blob,Type,Tag,Zettel,Config,InventoryList,Repo",
			inputs:            []string{"!md."},
		},
		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			defaultGenre:      ids.MakeGenre(genres.All()...),
			expectedOptimized: "-etikett-two.Zettel",
			expected:          "-etikett-two.Zettel",
			inputs:            []string{"-etikett-two.z"},
		},

		{
			stackInfo:         test_logz.MakeStackInfo(&t, 0),
			defaultGenre:      ids.MakeGenre(genres.All()...),
			expectedOptimized: "/repo:Repo",
			expected:          "/repo:Repo",
			inputs:            []string{"/repo:k"},
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
					tc.defaultGenre,
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
