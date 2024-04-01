package sku_fmt

import (
	"bytes"
	"io"
	"sort"
	"strings"
	"testing"

	"code.linenisgreat.com/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"github.com/DataDog/zstd"
)

func TestOne(t1 *testing.T) {
	t := test_logz.T{T: t1}

	b := new(bytes.Buffer)
	f := objekte_format.Default()
	o := objekte_format.Options{Tai: true}
	w := zstd.NewWriter(b)

	printer := MakeFormatBestandsaufnahmePrinter(w, f, o)

	sk := &sku.Transacted{}
	t.AssertNoError(sk.Kennung.SetWithKennung(kennung.MustHinweis("one/uno")))

	n, err := printer.Print(sk)

	{
		expected := int64(47)

		if n != expected {
			t.Errorf("expected %d but got %d", expected, n)
		}
	}

	{
		if err != nil {
			t.AssertNoError(err)
		}
	}

	sk = &sku.Transacted{}
	t.AssertNoError(sk.Kennung.SetWithKennung(kennung.MustHinweis("two/dos")))
	n, err = printer.Print(sk)

	{
		expected := int64(43)

		if n != expected {
			t.Errorf("expected %d but got %d", expected, n)
		}
	}

	{
		if err != nil {
			t.AssertNoError(err)
		}
	}

	w.Flush()

	op := objekte_format.Options{}

	scanner := MakeFormatBestandsaufnahmeScanner(zstd.NewReader(b), f, op)

	if !scanner.Scan() {
		t.Logf("scan error: %q", scanner.Error())
		t.Fatalf("expected ok scan")
	}

	_ = scanner.GetTransacted()
	t.AssertNoError(scanner.Error())

	if !scanner.Scan() {
		t.Fatalf("expected ok scan")
	}

	_ = scanner.GetTransacted()
	t.AssertNoError(scanner.Error())

	if scanner.Scan() {
		t.Fatalf("expected end scan")
	}

	t.AssertNoError(scanner.Error())
}

func TestBigMac(t1 *testing.T) {
	t := test_logz.T{T: t1}

	dataRaw := getRawData()

	dataComp := new(bytes.Buffer)

	comp := zstd.NewWriter(dataComp)

	io.Copy(comp, strings.NewReader(dataRaw))

	comp.Flush()

	f := objekte_format.Default()
	op := objekte_format.Options{}

	scanner := MakeFormatBestandsaufnahmeScanner(
		zstd.NewReader(dataComp),
		f,
		op,
	)

	i := 0

	for scanner.Scan() {
		sk := scanner.GetTransacted()

		if sk == nil {
			t.Errorf("expected sku but got nil")
		}

		i++

		t.AssertNoError(scanner.Error())
	}

	expected := 55

	if i != expected {
		t.Errorf("expected %d entries but got %d", expected, i)
	}
}

func TestOffsets(t1 *testing.T) {
	t := test_logz.T{T: t1}

	dataRaw := getRawData()

	f := objekte_format.Default()
	op := objekte_format.Options{Tai: true}

	scanner := MakeFormatBestandsaufnahmeScanner(
		strings.NewReader(dataRaw),
		f,
		op,
	)

	skus := make([]*sku.Transacted, 0)

	for scanner.Scan() {
		sk := scanner.GetTransacted()
		sk2 := &sku.Transacted{}
		t.AssertNoError(sk2.SetFromSkuLike(sk))
		skus = append(skus, sk2)
	}

	t.AssertNoError(scanner.Error())

	lookup := make(map[int64]*sku.Transacted)
	var b bytes.Buffer

	printer := MakeFormatBestandsaufnahmePrinter(&b, f, op)

	sk := &sku.Transacted{}
	t.AssertNoError(sk.Kennung.SetWithKennung(kennung.MustHinweis("one/uno")))

	for _, s := range skus {
		off := printer.Offset()
		_, err := printer.Print(s)
		t.AssertNoError(err)
		lookup[off] = s
	}

	bs := bytes.NewReader(b.Bytes())
	rb := catgut.MakeRingBuffer(bs, 0)

	sortedLookup := make([]int, 0, len(lookup))

	for off := range lookup {
		sortedLookup = append(sortedLookup, int(off))
	}

	sort.IntSlice(sortedLookup).Sort()

	for _, off := range sortedLookup {
		s := lookup[int64(off)]
		// t.Logf("at %d", off)
		_, err := rb.Seek(int64(off), io.SeekStart)
		if err != io.EOF {
			t.AssertNoError(err)
		}
		// re := rb.PeekReadable().String()

		sk = &sku.Transacted{}
		_, err = f.ParsePersistentMetadatei(rb, sk, op)
		// t.AssertErrorEquals(objekte_format.ErrV4ExpectedSpaceSeparatedKey, err)

		if !s.Equals(sk) {
			// t.Logf("\n%s", re)
			t.Errorf("expected %s but got %s", s, sk)
		}
	}
}

func getRawData() string {
	return `---
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ png
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ png
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ jpg
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2038907140.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2039724665.0
Typ pdf
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2039724665.0
Typ txt
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2039913141.0
Typ gpg
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2040580662.0
Typ pdf
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2041010005.0
Typ jpeg
---
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2041024286.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2041119818.0
Typ ttf
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2041120219.0
Typ pdf
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2041120264.0
Typ pdf
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2041761909.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2042394778.0
Typ m4r
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2043145169.0
Typ md
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2045916290.0
Typ toml-bookmark
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2049138420.0
Typ ttf
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2049673708.0
Typ png
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2054347577.86051
Typ pdf
---
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2070274012.190275
Typ task
---
Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Bezeichnung XXXX
Etikett XXXX
Etikett XXXX
Gattung Zettel
Kennung one/uno
Tai 2070624013.945645
Typ toml-bookmark
---
`
}
