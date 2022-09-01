package kennung

import "testing"

func TestEquals(t *testing.T) {
	p := Kennung{Left: 1, Right: 1}
	p1 := Kennung{Left: 1, Right: 1}

	if !p.Equals(p1) {
		t.Errorf("expected equality")
	}
}

func TestNotEquals(t *testing.T) {
	p := Kennung{Left: 1, Right: 1}
	p1 := Kennung{Left: 0, Right: 1}

	if p.Equals(p1) {
		t.Errorf("expected inequality")
	}
}

func TestToId1(t *testing.T) {
	assertToId(
		t,
		Kennung{Left: 0, Right: 0},
		1,
	)
}

func TestToId2(t *testing.T) {
	assertToId(
		t,
		Kennung{Left: 0, Right: 1},
		2,
	)
}

func TestToId42(t *testing.T) {
	assertToId(
		t,
		Kennung{Left: 5, Right: 3},
		42,
	)
}

func TestFromId5(t *testing.T) {
	assertFromId(t, "5", Kennung{Left: 1, Right: 1})
}

func TestFromId745(t *testing.T) {
	assertFromId(t, "745", Kennung{Left: 3, Right: 35})
}

func TestFromId10469(t *testing.T) {
	assertFromId(t, "10469", Kennung{Left: 28, Right: 116})
}

func TestFromId1(t *testing.T) {
	assertFromId(t, "1", Kennung{Left: 0, Right: 0})
}

func TestFromId2(t *testing.T) {
	assertFromId(t, "2", Kennung{Left: 0, Right: 1})
}

func TestFromId3(t *testing.T) {
	assertFromId(t, "3", Kennung{Left: 1, Right: 0})
}

func TestFromId42(t *testing.T) {
	assertFromId(t, "42", Kennung{Left: 5, Right: 3})
}

func TestFromId567(t *testing.T) {
	assertFromId(t, "567", Kennung{Left: 5, Right: 28})
}

func TestFromId672(t *testing.T) {
	assertFromId(t, "672", Kennung{Left: 5, Right: 31})
}

func assertFromId(t *testing.T, n string, expected Kennung) {
	t.Helper()

	p := &Kennung{}
	p.Set(n)

	if !p.Equals(expected) {
		t.Errorf("expected %v but got %v", expected, p)
	}
}

func assertToId(t *testing.T, p Kennung, expected Int) {
	t.Helper()

	id := p.Id()

	if id != expected {
		t.Errorf("expected %d but got %d", expected, id)
	}
}
