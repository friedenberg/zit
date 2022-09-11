package tridex

import (
	"os"
	"testing"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/test_logz"
)

type testStringer string

func TestMain(m *testing.M) {
	errors.SetTesting()
	code := m.Run()
	os.Exit(code)
}

func TestContains(t *testing.T) {
	sut := Make(
		"123456",
		"654321",
		"5",
	)

	expectedContains := []string{
		"123456",
		"654321",
		"5",
  }

	for _, e := range expectedContains {
		if !sut.Contains(e) {
			test_logz.Errorf(t, "expected %v to contain %s", sut, e)
		}
	}

	expectedNotContains := []string{
		"1",
		"12",
		"123",
		"1234",
		"12345",
		"1234567",
		"12X45",
		"3",
		"6",
		"65",
		"654",
		"6543",
		"65432",
	}

	for _, e := range expectedNotContains {
		if sut.Contains(e) {
			test_logz.Errorf(t, "expected %v to not contain %s", sut, e)
		}
	}
}

func TestAbbreviateOrphan(t *testing.T) {
	sut := Make(
		"one",
	)

	expectedContains := map[string]string{
		"one": "o",
	}

	for e, c := range expectedContains {
		if ca := sut.Abbreviate(e); ca != c {
			test_logz.Errorf(t, "%q: expected shorted length %q but got %q", e, c, ca)
		}
	}
}

func TestAbbreviateDegenerate(t *testing.T) {
	sut := Make(
		"mewtwo",
		"mew",
	)

	expectedContains := map[string]string{
		"mewtwo": "mewt",
		"mew":    "mew",
	}

	for e, c := range expectedContains {
		if ca := sut.Abbreviate(e); ca != c {
			test_logz.Printf("%#v", sut)
			test_logz.Errorf(t, "%q: expected shorted length %q but got %q", e, c, ca)
		}
	}
}

func TestExpandDegenerate(t *testing.T) {
	sut := Make(
		"mewtwo",
		"mew",
	)

	expectedContains := map[string]string{
		"mewt": "mewtwo",
		"mew":    "mew",
	}

	for e, c := range expectedContains {
		if ca := sut.Expand(e); ca != c {
			test_logz.Printf("%#v", sut)
			test_logz.Errorf(t, "%q: expected expanded %q but got %q", e, c, ca)
		}
	}
}

func TestAbbreviate(t *testing.T) {
	sut := Make(
		"12",
		"121",
		"127",
		"128",
		"123456",
		"654321",
		"mewtwo",
		"mew",
	)

	expectedContains := map[string]string{
		"mewtwo":   "mewt",
		"mew":      "mew",
		"121":      "121",
		"12":       "12",
		"123":      "123",
		"123456":   "123",
		"1234567":  "1234",
		"12345678": "1234",
		"124":      "124",
		"2":        "2",
	}

	for e, c := range expectedContains {
		if ca := sut.Abbreviate(e); ca != c {
			test_logz.Print(t, "%#v", sut)
			test_logz.Errorf(t, "%q: expected shorted length %q but got %q", e, c, ca)
		}
	}
}

func TestExpandOrphan(t *testing.T) {
	sut := Make(
		"654321",
	)

	expectedContains := map[string]string{
		"6":      "654321",
		"654321": "654321",
	}

	for a, e := range expectedContains {
		if ca := sut.Expand(a); ca != e {
			test_logz.Errorf(t, "%q: expected expanded %q but got %q", e, e, ca)
		}
	}
}

func TestExpand(t *testing.T) {
	sut := Make(
		"12",
		"121",
		"127",
		"128",
		"123456",
		"654321",
	)

	expectedContains := map[string]string{
		"6":      "654321",
		"654321": "654321",
		"128":    "128",
		"123":    "123456",
		"1232":   "",
	}

	for a, e := range expectedContains {
		if ca := sut.Expand(a); ca != e {
			test_logz.Errorf(t, "%q: expected expanded %q but got %q", e, e, ca)
		}
	}
}

// func TestShas(t *testing.T) {
// 	sut := Make(
// 		"7a2be8c643edd96b0cce2a1be32de30967a6db1f362047b954401458dd530f",
// 		"86324762e1fd27008f2b6c276ba82d3bbae349fd9819865a44931fae6bd41b",
// 		"905f0ed00a076771d48628477ff056e9b2e9d9ddcdafe4d78859517530f602",
// 		"a413bdb43426d528dc865d31e399719449b23481f5b9c10b8465a226fa863a",
// 		"b55eea6909f4b89967f531654dacb285723db6cd07d38f05de945f840db83d",
// 		"bc5784f13699736c709445114bf0374f9090c715549a826a2ad81a97ce055d",
// 		"be682a0034e409f23404b6588a7f39be6fe6eed333f341c9b97677ffa90303",
// 		"d123661f09fe5bd09d840b44c118c5c32cb8d0f33d1b9284aca95afdf18da9",
// 		"d45e1e12fececfc8613fa33bd237a5e511944234cb634ac05159b62438d15a",
// 		"e66bbf9a57566858d63522d83f4a3e17d7bdd8544c59d5c410be355b08b96a",
// 		"ec3f2668d283201e64db299bbdcc79e0885467172ea53246fef3c8880db14a",
// 		"ec5983b7f7d86cab238db64d0e8ea266dbbe3bf9133d4aa2acb61328abf94b",
// 		"f66f17833d028ca4aa311474a2e63f19ee6b46e379066defaf356f2a5405c8",
// 		"f8ff837c4d072cc7e28da224df391126952f6b6a6af237cdae00c663228c61",
// 		"ff05ec4214ad07c3f5f53ca6ed292a6115486f6864bb0612fa0e2d5fce7bec",
// 	)

// 	expectedContains := map[string]string{
// 		"6":    "654321",
// 		"128":  "128",
// 		"123":  "123456",
// 		"1232": "",
// 	}

// 	b, _ := json.Marshal(sut)

// 	for a, e := range expectedContains {
// 		if ca := sut.Expand(a); ca != e {
// 			test_logz.Errorf(t, "%#v: expected expanded %q but got %q", string(b), e, ca)
// 		}
// 	}
// }
