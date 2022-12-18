package zettel

import (
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/friedenberg/zit/src/bravo/test_logz"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/test_metadatei_io"
)

type noopCloser struct {
	*strings.Reader
}

func (c noopCloser) Close() error {
	return nil
}

type akteReaderFactory struct {
	t     test_logz.T
	akten map[string]string
}

func (arf akteReaderFactory) AkteReader(s sha.Sha) (r sha.ReadCloser, err error) {
	var v string
	var ok bool

	if v, ok = arf.akten[s.String()]; !ok {
		arf.t.Fatalf("request for non-existent akte: %s", s)
	}

	r = sha.MakeNopReadCloser(io.NopCloser(strings.NewReader(v)))

	return
}

func writeFormat(t test_logz.T, z Objekte, f Format, includeAkte bool, akteBody string) (out string) {
	hash := sha256.New()
	_, err := io.Copy(hash, strings.NewReader(akteBody))

	if err != nil {
		t.Fatalf("%s", err)
	}

	akteShaRaw := fmt.Sprintf("%x", hash.Sum(nil))
	var akteSha sha.Sha

	if err := akteSha.Set(akteShaRaw); err != nil {
		t.Fatalf("%s", err)
	}

	z.Akte = akteSha

	sb := &strings.Builder{}

	c := FormatContextWrite{
		Zettel:      z,
		Out:         sb,
		IncludeAkte: includeAkte,
	}

	if _, err := f.WriteTo(c); err != nil {
		t.Errorf("%s", err)
	}

	out = sb.String()

	return
}

func TestWriteWithoutAkte(t1 *testing.T) {
	t := test_logz.T{T: t1}

	z := Objekte{
		Bezeichnung: bezeichnung.Make("the title"),
		Etiketten: makeEtiketten(t,
			"tag1",
			"tag2",
			"tag3",
		),
		Typ: makeAkteExt(t, "md"),
	}

	actual := writeFormat(t, z, textParser{}, false, `the body`)

	expected := `---
# the title
- tag1
- tag2
- tag3
! fa8242e99f48966ca514092b4233b446851f42b57ad5031bf133e1dd76787f3e.md
---
`

	if expected != actual {
		t.Fatalf("zettel:\nexpected: %#v\n  actual: %#v", expected, actual)
	}
}

func TestWriteWithInlineAkte(t1 *testing.T) {
	t := test_logz.T{T: t1}

	z := Objekte{
		Bezeichnung: bezeichnung.Make("the title"),
		Etiketten: makeEtiketten(t,
			"tag1",
			"tag2",
			"tag3",
		),
		Typ: makeAkteExt(t, "md"),
	}

	format := textParser{
		AkteFactory: test_metadatei_io.FixtureFactoryReadWriteCloser(
			map[string]string{
				"fa8242e99f48966ca514092b4233b446851f42b57ad5031bf133e1dd76787f3e": "the body\n",
			},
		),
	}

	actual := writeFormat(t, z, format, true, `the body`)

	expected := `---
# the title
- tag1
- tag2
- tag3
! md
---

the body
`

	if expected != actual {
		t.Fatalf("zettel:\nexpected: %#v\n  actual: %#v", expected, actual)
	}
}
