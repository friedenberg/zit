package trie

import "testing"

type testStringer string

func (ts testStringer) String() string {
	return string(ts)
}

func TestContains(t *testing.T) {
	sut := Make(
		testStringer("123456"),
		testStringer("654321"),
	)

	expectedContains := []string{
		"1",
		"12",
		"123",
		"1234",
		"12345",
		"123456",
		"654321",
		"65432",
		"6543",
		"654",
		"65",
		"6",
	}

	for _, e := range expectedContains {
		es := testStringer(e)

		if !sut.Contains(es) {
			t.Errorf("expected %v to contain %s", sut, es)
		}
	}

	expectedNotContains := []string{
		"3",
		"12X45",
		"1234567",
	}

	for _, e := range expectedNotContains {
		es := testStringer(e)

		if sut.Contains(es) {
			t.Errorf("expected %v to not contain %s", sut, es)
		}
	}
}

func TestShortestUnique(t *testing.T) {
	sut := Make(
		testStringer("12"),
		testStringer("121"),
		testStringer("127"),
		testStringer("128"),
		testStringer("123456"),
		testStringer("654321"),
	)

	expectedContains := map[string]string{
		"123":      "123",
		"123456":   "123",
		"1234567":  "1234567",
		"12345678": "1234567",
		"124":      "124",
		"2":        "2",
	}

	for e, c := range expectedContains {
		es := testStringer(e)

		if ca := sut.Abbreviate(es); ca != c {
			t.Errorf("%q: expected shorted length %q but got %q", es, c, ca)
		}
	}
}
