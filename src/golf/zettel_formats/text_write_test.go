package zettel_formats

import (
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
)

type noopCloser struct {
	*strings.Reader
}

func (c noopCloser) Close() error {
	return nil
}

type akteReaderFactory struct {
	t     *testing.T
	akten map[string]string
}

func (arf akteReaderFactory) AkteReader(s sha.Sha) (r io.ReadCloser, err error) {
	var v string
	var ok bool

	if v, ok = arf.akten[s.String()]; !ok {
		arf.t.Fatalf("request for non-existent akte: %s", s)
	}

	r = noopCloser{strings.NewReader(v)}

	return
}

func writeFormat(t *testing.T, z zettel.Zettel, f zettel.Format, includeAkte bool, akteBody string) (out string) {
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

	c := zettel.FormatContextWrite{
		Zettel:      z,
		Out:         sb,
		IncludeAkte: includeAkte,
		AkteReaderFactory: akteReaderFactory{
			t: t,
			akten: map[string]string{
				akteShaRaw: akteBody,
			},
		},
	}

	if _, err := f.WriteTo(c); err != nil {
		t.Errorf("%s", err)
	}

	out = sb.String()

	return
}

func TestWriteWithoutAkte(t *testing.T) {
	z := zettel.Zettel{
		Bezeichnung: "the title",
		Etiketten: makeEtiketten(t,
			"tag1",
			"tag2",
			"tag3",
		),
		AkteExt: makeAkteExt(t, "md"),
	}

	actual := writeFormat(t, z, Text{}, false, `the body`)

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

func TestWriteWithInlineAkte(t *testing.T) {
	z := zettel.Zettel{
		Bezeichnung: "the title",
		Etiketten: makeEtiketten(t,
			"tag1",
			"tag2",
			"tag3",
		),
		AkteExt: makeAkteExt(t, "md"),
	}

	actual := writeFormat(t, z, Text{}, true, `the body`)

	expected := `---
# the title
- tag1
- tag2
- tag3
! md
---

the body`

	if expected != actual {
		t.Fatalf("zettel:\nexpected: %#v\n  actual: %#v", expected, actual)
	}
}
