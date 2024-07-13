package store_verzeichnisse

import (
	"bytes"
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func TestBinaryOne(t1 *testing.T) {
	t := test_logz.T{T: t1}

	b := new(bytes.Buffer)
	coder := binaryEncoder{Sigil: kennung.SigilSchwanzen}
	decoder := makeBinary(kennung.SigilSchwanzen)
	expected := &sku.Transacted{}
	var expectedN int64
	var err error

	{
		t.AssertNoError(expected.Kennung.SetWithIdLike(kennung.MustHinweis("one/uno")))
		expected.SetTai(kennung.NowTai())
		t.AssertNoError(expected.Metadatei.Akte.Set(
			"ed500e315f33358824203cee073893311e0a80d77989dc55c5d86247d95b2403",
		))
		t.AssertNoError(expected.Metadatei.Typ.Set("da-typ"))
		t.AssertNoError(expected.Metadatei.Bezeichnung.Set("the bez"))
		t.AssertNoError(expected.AddEtikettPtr(kennung.MustEtikettPtr("tag")))
		t.AssertNoError(expected.Metadatei.Mutter().Set(
			"3c5d8b1db2149d279f4d4a6cb9457804aac6944834b62aa283beef99bccd10f0",
		))
		t.AssertNoError(expected.CalculateObjekteShas())

		t.Logf("%s", expected)

		expectedN, err = coder.writeFormat(b, skuWithSigil{Transacted: expected})
		t.AssertNoError(err)
	}

	actual := skuWithRangeAndSigil{
		skuWithSigil: skuWithSigil{
			Transacted: &sku.Transacted{},
		},
	}

	{
		n, err := decoder.readFormatAndMatchSigil(b, &actual)
		t.AssertNoError(err)
		t.Logf("%s", actual)

		{
			if n != expectedN {
				t.Errorf("expected %d but got %d", expectedN, n)
			}
		}
	}

	if !sku.TransactedEqualer.Equals(expected, actual.Transacted) {
		t.NotEqual(expected, actual)
	}
}

// func TestBigMac(t1 *testing.T) {
// 	t := test_logz.T{T: t1}

// 	dataRaw := getRawData()

// 	dataComp := new(bytes.Buffer)

// 	comp := zstd.NewWriter(dataComp)

// 	io.Copy(comp, strings.NewReader(dataRaw))

// 	comp.Flush()

// 	f := objekte_format.Default()
// 	op := objekte_format.Options{}

// 	scanner := MakeFormatBestandsaufnahmeScanner(
// 		zstd.NewReader(dataComp),
// 		f,
// 		op,
// 	)

// 	i := 0

// 	for scanner.Scan() {
// 		sk := scanner.GetTransacted()

// 		if sk == nil {
// 			t.Errorf("expected sku but got nil")
// 		}

// 		i++

// 		t.AssertNoError(scanner.Error())
// 	}

// 	expected := 55

// 	if i != expected {
// 		t.Errorf("expected %d entries but got %d", expected, i)
// 	}
// }

// func TestOffsets(t1 *testing.T) {
// 	t := test_logz.T{T: t1}

// 	dataRaw := getRawData()

// 	f := objekte_format.Default()
// 	op := objekte_format.Options{Tai: true}

// 	scanner := MakeFormatBestandsaufnahmeScanner(
// 		strings.NewReader(dataRaw),
// 		f,
// 		op,
// 	)

// 	skus := make([]*sku.Transacted, 0)

// 	for scanner.Scan() {
// 		sk := scanner.GetTransacted()
// 		sk2 := &sku.Transacted{}
// 		t.AssertNoError(sk2.SetFromSkuLike(sk))
// 		skus = append(skus, sk2)
// 	}

// 	t.AssertNoError(scanner.Error())

// 	lookup := make(map[int64]*sku.Transacted)
// 	var b bytes.Buffer

// 	printer := MakeFormatBestandsaufnahmePrinter(&b, f, op)

// 	sk := &sku.Transacted{}
// 	t.AssertNoError(sk.Kennung.SetWithKennung(kennung.MustHinweis("one/uno")))

// 	for _, s := range skus {
// 		off := printer.Offset()
// 		_, err := printer.Print(s)
// 		t.AssertNoError(err)
// 		lookup[off] = s
// 	}

// 	bs := bytes.NewReader(b.Bytes())
// 	rb := catgut.MakeRingBuffer(bs, 0)

// 	sortedLookup := make([]int, 0, len(lookup))

// 	for off := range lookup {
// 		sortedLookup = append(sortedLookup, int(off))
// 	}

// 	sort.IntSlice(sortedLookup).Sort()

// 	for _, off := range sortedLookup {
// 		s := lookup[int64(off)]
// 		t.Logf("at %d", off)
// 		_, err := rb.Seek(int64(off), io.SeekStart)
// 		if err != io.EOF {
// 			t.AssertNoError(err)
// 		}
// 		// re := rb.PeekReadable().String()

// 		sk = &sku.Transacted{}
// 		_, err = f.ParsePersistentMetadatei(rb, sk, op)
// 		// t.AssertErrorEquals(objekte_format.ErrV4ExpectedSpaceSeparatedKey, err)

// 		if !s.Equals(sk) {
// 			// t.Logf("\n%s", re)
// 			t.Errorf("expected %s but got %s", s, sk)
// 		}
// 	}
// }

// func getRawData() string {
// 	return `---
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ png
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ png
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ jpg
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Bezeichnung XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2038907140.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2039724665.0
// Typ pdf
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2039724665.0
// Typ txt
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2039913141.0
// Typ gpg
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2040580662.0
// Typ pdf
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2041010005.0
// Typ jpeg
// ---
// Bezeichnung XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2041024286.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2041119818.0
// Typ ttf
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2041120219.0
// Typ pdf
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2041120264.0
// Typ pdf
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2041761909.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2042394778.0
// Typ m4r
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2043145169.0
// Typ md
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2045916290.0
// Typ toml-bookmark
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2049138420.0
// Typ ttf
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2049673708.0
// Typ png
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2054347577.86051
// Typ pdf
// ---
// Bezeichnung XXXX
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2070274012.190275
// Typ task
// ---
// Akte e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
// Bezeichnung XXXX
// Etikett XXXX
// Etikett XXXX
// Gattung Zettel
// Kennung one/uno
// Tai 2070624013.945645
// Typ toml-bookmark
// ---
// `
// }
