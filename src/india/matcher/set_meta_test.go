package matcher

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/friedenberg/zit/src/bravo/test_logz"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/kennung"
)

func TestMetaSetGob(t1 *testing.T) {
	t := test_logz.T{T: t1}

	bs := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(bs)
	gob.Register(&metaSet{})

	{
		sut := MakeMetaSet(
			nil,
			kennung.Abbr{},
			nil,
			nil,
			gattungen.MakeSet(
				gattung.Zettel,
			),
			nil,
			kennung.Index{},
		)

		if err := sut.Set("one/uno@zettel"); err != nil {
			t.Errorf("setting failed: %s", err)
		}

		if err := enc.Encode(&sut); err != nil {
			t.Errorf("encoding failed: %s", err)
		}
	}

	dec := gob.NewDecoder(bs)

	{
		var sut MetaSet

		if err := dec.Decode(&sut); err != nil {
			t.Errorf("decoding failed: %s", err)
		}

		if sut == nil {
			t.Errorf("sut was nil")
		}
	}
}
