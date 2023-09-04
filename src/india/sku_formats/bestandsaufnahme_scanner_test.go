package sku_formats

import (
	"bytes"
	"io"
	"testing"

	"github.com/DataDog/zstd"
	"github.com/friedenberg/zit/src/bravo/test_logz"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/transacted"
)

func TestOne(t1 *testing.T) {
	t := test_logz.T{T: t1}

	b := new(bytes.Buffer)
	f := objekte_format.BestandsaufnahmeFormatIncludeTai()
	w := zstd.NewWriter(b)

	printer := MakeFormatBestandsaufnahmePrinter(w, f)

	n, err := printer.Print(transacted.Zettel{
		Kennung: kennung.MustHinweis("one/uno"),
	})

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

	n, err = printer.Print(transacted.Zettel{
		Kennung: kennung.MustHinweis("two/dos"),
	})

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

	scanner := MakeFormatbestandsaufnahmeScanner(zstd.NewReader(b), f)
	var sk sku.SkuLike

	sk, n, err = scanner.Scan()

	{
		if err != nil {
			t.AssertNoError(err)
		}
	}

	{
		sk1, ok := sk.(*transacted.Zettel)

		if !ok {
			t.Errorf("expected %T but got %T", sk1, sk)
		}
	}

	{
		expected := int64(47)

		if n != expected {
			t.Errorf("expected %d but got %d", expected, n)
		}
	}

	sk, n, err = scanner.Scan()

	{
		if err != nil {
			t.AssertNoError(err)
		}
	}

	{
		sk1, ok := sk.(*transacted.Zettel)

		if !ok {
			t.Errorf("expected %T but got %T", sk1, sk)
		}
	}

	{
		expected := int64(43)

		if n != expected {
			t.Errorf("expected %d but got %d", expected, n)
		}
	}

	sk, n, err = scanner.Scan()

	{
		if err != nil {
			t.AssertError(io.EOF)
		}
	}

	{
		expected := int64(0)

		if n != expected {
			t.Errorf("expected %d but got %d", expected, n)
		}
	}

	if sk != nil {
		t.Errorf("expected nil but got %s", sk)
	}
}
