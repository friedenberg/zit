package collections

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/friedenberg/zit/src/bravo/test_logz"
)

func TestBitset0CapGreaterAdd(t1 *testing.T) {
	t := test_logz.T{T: t1}

	sut := MakeBitset(20)
	sut.Add(19)

	if !sut.Get(19) {
		t.Errorf("expected bitset to contain idx %d", 19)
	}
}

func TestBitset1CapLessAdd(t1 *testing.T) {
	t := test_logz.T{T: t1}

	sut := MakeBitset(20)
	toAdd := int(21)
	sut.Add(toAdd)

	if !sut.Get(toAdd) {
		t.Errorf("expected bitset to contain idx %d", toAdd)
	}
}

func TestBitset2CapLessAddRemove(t1 *testing.T) {
	t := test_logz.T{T: t1}

	sut := MakeBitset(20)
	toAdd := int(256)
	sut.Add(toAdd)

	if !sut.Get(toAdd) {
		t.Errorf("expected bitset to contain idx %d", toAdd)
	}

	sut.Del(toAdd)

	if sut.Get(toAdd) {
		t.Errorf("expected bitset to not contain idx %d", toAdd)
	}
}

func TestBitset3WouldGrowTooLarge(t1 *testing.T) {
	t := test_logz.T{T: t1}

	defer func() {
		e := recover()

		if e == nil {
			t.Errorf("expected bitset to panic")
		}
	}()

	sut := MakeBitset(20)
	toAdd := int(MaxBitsetIdx + 1)
	sut.Add(toAdd)
}

func TestBitset4Gob(t1 *testing.T) {
	t := test_logz.T{T: t1}

	sut := MakeBitset(20)
	toAdd := 12
	sut.Add(toAdd)

	bytes := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(bytes)

	if err := enc.Encode(sut); err != nil {
		t.Errorf("expected gob.Encode to succeed: %s", err)
	}

	sut2 := MakeBitset(20)
	dec := gob.NewDecoder(bytes)

	if err := dec.Decode(sut2); err != nil {
		t.Errorf("expected gob.Decode to succeed: %s", err)
	}

	if !sut.Equals(sut2) {
		t.Errorf("expected equality")
	}
}

func TestBitset5Equals(t1 *testing.T) {
	t := test_logz.T{T: t1}

	sut := MakeBitset(20)
	toAdd := 12
	sut.Add(toAdd)

	sut2 := MakeBitset(20)
	sut2.Add(toAdd)

	if !sut.Equals(sut2) {
		t.Errorf("expected equality")
	}
}

func BenchmarkAdd(b *testing.B) {
	sut := MakeBitset(int(b.N))

	b.ResetTimer()

	j := int(0)

	for i := 0; i < b.N; i++ {
		if j > MaxBitsetIdx {
			j = 0
		}

		sut.Add(int(j))
		j++
	}
}
