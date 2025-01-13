package sku

import (
	"bytes"
	"encoding/gob"
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
)

func TestGob(t1 *testing.T) {
	t := test_logz.T{T: t1}

	type testType = Transacted

	var expected testType

	if err := expected.ObjectId.Set("test-tag"); err != nil {
		t.Fatalf("failed to set object id: %w", err)
	}

	var b bytes.Buffer

	enc := gob.NewEncoder(&b)

	if err := enc.Encode(&expected); err != nil {
		t.Fatalf("failed to encode config: %w", err)
	}

	dec := gob.NewDecoder(&b)

	var actual testType

	if err := dec.Decode(&actual); err != nil {
		t.Fatalf("failed to decode config: %w", err)
	}

	t.AssertNotEqual(expected.ObjectId.String(), actual.ObjectId.String())
}
