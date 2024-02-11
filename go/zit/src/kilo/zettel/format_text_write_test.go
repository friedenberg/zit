package zettel

import (
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
	"testing"

	"code.linenisgreat.com/zit-go/src/bravo/test_logz"
	"code.linenisgreat.com/zit-go/src/charlie/sha"
	"code.linenisgreat.com/zit-go/src/echo/bezeichnung"
	"code.linenisgreat.com/zit-go/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit-go/src/foxtrot/test_metadatei_io"
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

func writeFormat(
	t test_logz.T,
	m *metadatei.Metadatei,
	f metadatei.TextFormatter,
	includeAkte bool,
	akteBody string,
) (out string) {
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

	if err = m.Akte.SetShaLike(&akteSha); err != nil {
		t.Fatalf("%s", err)
	}

	sb := &strings.Builder{}

	if _, err := f.FormatMetadatei(sb, m); err != nil {
		t.Errorf("%s", err)
	}

	out = sb.String()

	return
}

func TestWriteWithoutAkte(t1 *testing.T) {
	t := test_logz.T{T: t1}

	z := &metadatei.Metadatei{
		Bezeichnung: bezeichnung.Make("the title"),
		Typ:         makeAkteExt(t, "md"),
	}

	z.SetEtiketten(makeEtiketten(t,
		"tag1",
		"tag2",
		"tag3",
	))

	af := test_metadatei_io.FixtureFactoryReadWriteCloser(
		map[string]string{
			"fa8242e99f48966ca514092b4233b446851f42b57ad5031bf133e1dd76787f3e": "the body\n",
		},
	)

	format := metadatei.MakeTextFormatterMetadateiOnly(
		af,
		nil,
	)

	actual := writeFormat(t, z, format, false, `the body`)

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

	z := &metadatei.Metadatei{
		Bezeichnung: bezeichnung.Make("the title"),
		Typ:         makeAkteExt(t, "md"),
	}

	z.SetEtiketten(makeEtiketten(t,
		"tag1",
		"tag2",
		"tag3",
	))

	af := test_metadatei_io.FixtureFactoryReadWriteCloser(
		map[string]string{
			"fa8242e99f48966ca514092b4233b446851f42b57ad5031bf133e1dd76787f3e": "the body\n",
		},
	)

	format := metadatei.MakeTextFormatterMetadateiInlineAkte(
		af,
		nil,
	)

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
