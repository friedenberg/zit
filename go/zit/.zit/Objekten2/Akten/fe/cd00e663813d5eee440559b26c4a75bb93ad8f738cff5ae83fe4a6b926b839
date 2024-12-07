package config

import (
	"bytes"
	"encoding/gob"
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func TestGob(t1 *testing.T) {
	t := test_logz.T{T: t1}

	ta := sku.GetTransactedPool().Get()

	if err := ta.ObjectId.Set("test-tag"); err != nil {
		t.Fatalf("failed to set object id: %w", err)
	}

	var b bytes.Buffer

	enc := gob.NewEncoder(&b)

	if err := enc.Encode(ta); err != nil {
		t.Fatalf("failed to encode config: %w", err)
	}

	dec := gob.NewDecoder(&b)

	var actual sku.Transacted

	if err := dec.Decode(&actual); err != nil {
		t.Fatalf("failed to decode config: %w", err)
	}

	t.AssertNotEqual(ta.ObjectId.String(), actual.ObjectId.String())
}
