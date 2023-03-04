package kennung

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/test_logz"
	"github.com/friedenberg/zit/src/delta/gattungen"
)

func TestMetaSetGob(t1 *testing.T) {
	t := test_logz.T{T: t1}

	bs := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(bs)
	gob.Register(&metaSet{})

	{
		sut := MakeMetaSet(
			Expanders{},
			gattungen.MakeSet(
				gattung.Zettel,
			),
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
