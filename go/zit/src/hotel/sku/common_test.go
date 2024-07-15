package sku

import (
	"crypto/sha256"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_ptr"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/test_object_metadata_io"
)

type inlineTypChecker struct {
	answer bool
}

func (t inlineTypChecker) IsInlineTyp(k ids.Type) bool {
	return t.answer
}

func makeTagSet(t test_logz.T, vs ...string) (es ids.TagSet) {
	var err error

	if es, err = collections_ptr.MakeValueSetString[ids.Tag, *ids.Tag](nil, vs...); err != nil {
		t.Fatalf("%s", err)
	}

	return
}

func makeBlobExt(t test_logz.T, v string) (es ids.Type) {
	if err := es.Set(v); err != nil {
		t.Fatalf("%s", err)
	}

	return
}

func readFormat(
	t1 test_logz.T,
	f object_metadata.TextFormat,
	af *test_object_metadata_io.BlobIOFactory,
	contents string,
) (z *object_metadata.Metadata, a string) {
	var zt Transacted

	t := t1.Skip(1)

	n, err := f.ParseMetadata(
		strings.NewReader(contents),
		&zt,
	)
	if err != nil {
		t.Fatalf("failed to read zettel format: %s", err)
	}

	if n != int64(len(contents)) {
		t.Fatalf("expected to read %d but only read %d", len(contents), n)
	}

	z = zt.GetMetadata()
	a = af.CurrentBufferString()

	return
}

func TestMakeTags(t1 *testing.T) {
	t := test_logz.T{T: t1}

	vs := []string{
		"tag1",
		"tag2",
		"tag3",
	}

	var sut ids.TagSet
	var err error

	if sut, err = ids.MakeTagSetStrings(vs...); err != nil {
		t.Fatalf("%s", err)
	}

	if sut.Len() != 3 {
		t.Fatalf("expected len 3 but got %d", sut.Len())
	}

	{
		ac := sut.Len()

		if ac != 3 {
			t.Fatalf("expected len 3 but got %d", ac)
		}
	}

	sut2 := sut.CloneSetLike()

	if sut2.Len() != 3 {
		t.Fatalf("expected len 3 but got %d", sut2.Len())
	}

	{
		ac := iter.SortedStrings[ids.Tag](sut)

		if !reflect.DeepEqual(ac, vs) {
			t.Fatalf("expected %q but got %q", vs, ac)
		}
	}

	{
		ex := "tag1, tag2, tag3"
		ac := iter.StringCommaSeparated[ids.Tag](sut)

		if ac != ex {
			t.Fatalf("expected %q but got %q", ex, ac)
		}
	}

	{
		ex := "tag1, tag2, tag3"
		ac := iter.StringCommaSeparated[ids.Tag](
			sut.CloneSetLike(),
		)

		if ac != ex {
			t.Fatalf("expected %q but got %q", ex, ac)
		}
	}
}

func TestEqualitySelf(t1 *testing.T) {
	t := test_logz.T{T: t1}

	text := &object_metadata.Metadata{
		Description: descriptions.Make("the title"),
		Type:        makeBlobExt(t, "text"),
	}

	text.SetTags(makeTagSet(t,
		"tag1",
		"tag2",
		"tag3",
	))

	if !text.Equals(text) {
		t.Fatalf("expected %v to equal itself", text)
	}
}

func TestEqualityNotSelf(t1 *testing.T) {
	t := test_logz.T{T: t1}

	text := object_metadata.Metadata{
		Description: descriptions.Make("the title"),
		Type:        makeBlobExt(t, "text"),
	}

	text.SetTags(makeTagSet(t,
		"tag1",
		"tag2",
		"tag3",
	))

	text1 := &object_metadata.Metadata{
		Description: descriptions.Make("the title"),
		Type:        makeBlobExt(t, "text"),
	}

	text1.SetTags(makeTagSet(t,
		"tag1",
		"tag2",
		"tag3",
	))

	if !text.Equals(text1) {
		t.Fatalf("expected %v to equal %v", text, text1)
	}
}

func makeTestTextFormat(
	af *test_object_metadata_io.BlobIOFactory,
) object_metadata.TextFormat {
	if af == nil {
		af = test_object_metadata_io.FixtureFactoryReadWriteCloser(nil)
	}

	return object_metadata.MakeTextFormat(
		af,
		nil,
	)
}

func TestReadWithoutBlob(t1 *testing.T) {
	t := test_logz.T{T: t1}
	af := test_object_metadata_io.FixtureFactoryReadWriteCloser(nil)

	actual, blob := readFormat(
		t,
		makeTestTextFormat(af),
		af,
		`---
# the title
- tag1
- tag2
- tag3
! md
---
`,
	)

	expected := &object_metadata.Metadata{
		Description: descriptions.Make("the title"),
		Type:        makeBlobExt(t, "md"),
	}

	expected.SetTags(makeTagSet(t,
		"tag1",
		"tag2",
		"tag3",
	))

	if !actual.Equals(expected) {
		t.Fatalf("zettel:\nexpected: %#v\n  actual: %#v", expected, actual)
	}

	if blob != "" {
		t.Fatalf("blob:\nexpected empty but got %q", blob)
	}
}

func TestReadWithoutBlobWithMultilineDescription(t1 *testing.T) {
	t := test_logz.T{T: t1}

	af := test_object_metadata_io.FixtureFactoryReadWriteCloser(nil)

	actual, blob := readFormat(
		t,
		makeTestTextFormat(af),
		af,
		`---
# the title
# continues
- tag1
- tag2
- tag3
! md
---
`,
	)

	expected := &object_metadata.Metadata{
		Description: descriptions.Make("the title\ncontinues"),
		Type:        makeBlobExt(t, "md"),
	}

	expected.SetTags(makeTagSet(t,
		"tag1",
		"tag2",
		"tag3",
	))

	if !actual.Equals(expected) {
		t.Fatalf("zettel:\nexpected: %#v\n  actual: %#v", expected, actual)
	}

	if blob != "" {
		t.Fatalf("blob:\nexpected empty but got %q", blob)
	}
}

func TestReadWithBlob(t1 *testing.T) {
	t := test_logz.T{T: t1}

	af := test_object_metadata_io.FixtureFactoryReadWriteCloser(
		map[string]string{
			"fa8242e99f48966ca514092b4233b446851f42b57ad5031bf133e1dd76787f3e": "the body\n",
			"036a8e44e472523c0306946f2712f372c234f8a24532e933f1509ae4db0da064": "the body",
		},
	)

	actual, blob := readFormat(
		t,
		makeTestTextFormat(af),
		af,
		`---
# the title
- tag1
- tag2
- tag3
! md
---

the body
`,
	)

	expected := &object_metadata.Metadata{
		Description: descriptions.Make("the title"),
		Type:        makeBlobExt(t, "md"),
	}

	errors.PanicIfError(expected.Blob.Set(
		"036a8e44e472523c0306946f2712f372c234f8a24532e933f1509ae4db0da064",
	))

	expected.SetTags(makeTagSet(t,
		"tag1",
		"tag2",
		"tag3",
	))

	if !actual.Equals(expected) {
		t.Fatalf("zettel:\nexpected: %#v\n  actual: %#v", expected, actual)
	}

	expectedBlob := "the body\n"

	if expectedBlob != blob {
		t.Fatalf("blob:\nexpected: %#v\n  actual: %#v", expectedBlob, blob)
	}
}

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
		Type:        makeBlobExt(t, "md"),
	}

	z.SetTags(makeTagSet(t,
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
		Type:        makeBlobExt(t, "md"),
	}

	z.SetTags(makeTagSet(t,
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
