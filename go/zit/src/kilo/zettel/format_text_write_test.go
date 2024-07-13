package zettel

import (
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/test_object_metadata_io"
)

type noopCloser struct {
	*strings.Reader
}

func (c noopCloser) Close() error {
	return nil
}

type blobReaderFactory struct {
	t     test_logz.T
	blobs map[string]string
}

func (arf blobReaderFactory) BlobReader(s sha.Sha) (r sha.ReadCloser, err error) {
	var v string
	var ok bool

	if v, ok = arf.blobs[s.String()]; !ok {
		arf.t.Fatalf("request for non-existent blob: %s", s)
	}

	r = sha.MakeNopReadCloser(io.NopCloser(strings.NewReader(v)))

	return
}

func writeFormat(
	t test_logz.T,
	m *object_metadata.Metadata,
	f object_metadata.TextFormatter,
	includeBlob bool,
	blobBody string,
) (out string) {
	hash := sha256.New()
	_, err := io.Copy(hash, strings.NewReader(blobBody))
	if err != nil {
		t.Fatalf("%s", err)
	}

	blobShaRaw := fmt.Sprintf("%x", hash.Sum(nil))
	var blobSha sha.Sha

	if err := blobSha.Set(blobShaRaw); err != nil {
		t.Fatalf("%s", err)
	}

	if err = m.Blob.SetShaLike(&blobSha); err != nil {
		t.Fatalf("%s", err)
	}

	sb := &strings.Builder{}

	if _, err := f.FormatMetadata(sb, m); err != nil {
		t.Errorf("%s", err)
	}

	out = sb.String()

	return
}

func TestWriteWithoutBlob(t1 *testing.T) {
	t := test_logz.T{T: t1}

	z := &object_metadata.Metadata{
		Description: descriptions.Make("the title"),
		Type:         makeBlobExt(t, "md"),
	}

	z.SetEtiketten(makeEtiketten(t,
		"tag1",
		"tag2",
		"tag3",
	))

	af := test_object_metadata_io.FixtureFactoryReadWriteCloser(
		map[string]string{
			"fa8242e99f48966ca514092b4233b446851f42b57ad5031bf133e1dd76787f3e": "the body\n",
		},
	)

	format := object_metadata.MakeTextFormatterMetadataOnly(
		object_metadata.TextFormatterOptions{},
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

func TestWriteWithInlineBlob(t1 *testing.T) {
	t := test_logz.T{T: t1}

	z := &object_metadata.Metadata{
		Description: descriptions.Make("the title"),
		Type:         makeBlobExt(t, "md"),
	}

	z.SetEtiketten(makeEtiketten(t,
		"tag1",
		"tag2",
		"tag3",
	))

	af := test_object_metadata_io.FixtureFactoryReadWriteCloser(
		map[string]string{
			"fa8242e99f48966ca514092b4233b446851f42b57ad5031bf133e1dd76787f3e": "the body\n",
		},
	)

	format := object_metadata.MakeTextFormatterMetadataInlineBlob(
		object_metadata.TextFormatterOptions{},
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
