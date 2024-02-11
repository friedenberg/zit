package kennung

import (
	"bytes"
	"encoding/gob"
	"testing"

	"code.linenisgreat.com/zit/src/bravo/test_logz"
)

type BobTest interface {
	GetValue() string
}

type bobTest struct {
	v string
}

func (s bobTest) GetValue() string {
	return s.v
}

func (s bobTest) MarshalBinary() (bs []byte, err error) {
	bs = []byte(s.v)

	return
}

func (s *bobTest) UnmarshalBinary(bs []byte) (err error) {
	s.v = string(bs)

	return
}

func TestBob(t1 *testing.T) {
	t := test_logz.T{T: t1}

	bs := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(bs)
	gob.Register(&bobTest{})

	{
		var sut BobTest

		sut = &bobTest{v: "wow"}

		if err := enc.Encode(&sut); err != nil {
			t.Errorf("encoding failed: %s", err)
		}
	}

	dec := gob.NewDecoder(bs)

	{
		var sut BobTest

		if err := dec.Decode(&sut); err != nil {
			t.Errorf("decoding failed: %s", err)
		}

		if sut == nil {
			t.Errorf("sut was nil")
		}

		if sut.GetValue() != "wow" {
			t.Errorf("expected value %s but got %s", "wow", sut.GetValue())
		}
	}
}
