package tag_paths

import (
	"testing"

	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
)

func TestAddPaths(t1 *testing.T) {
	t := test_logz.T{T: t1}

	var es Etiketten

	es.AddPath(MakePathWithType(
		catgut.MakeFromString("area-home"),
		catgut.MakeFromString("project-reno"),
	))

	{
		i, ok := es.All.ContainsEtikett(catgut.MakeFromString("area"))

		if !ok {
			t.Errorf("expected some etikett: %d, %t, %s", i, ok, es)
		}
	}

	es.AddPath(MakePathWithType(
		catgut.MakeFromString("area-career"),
		catgut.MakeFromString("project-recurse"),
	))

	{
		i, ok := es.All.ContainsEtikett(catgut.MakeFromString("area"))

		if !ok {
			t.Errorf("expected some etikett: %d, %t, %s", i, ok, es.All)
		}
	}
}

func TestRealWorld(t1 *testing.T) {
	t := test_logz.T{T: t1}

	var es Etiketten

	es.AddEtikett(catgut.MakeFromString("pom-1"))
	es.AddEtikett(catgut.MakeFromString("req-comp-internet"))
	es.AddEtikett(catgut.MakeFromString("today-in_progress"))

	{
		e := catgut.MakeFromString("req-comp-internet")
		_, ok := es.All.ContainsEtikett(e)

		if !ok {
			t.Errorf("expected %s to be in %s", e, es)
		}
	}

	es.AddPath(MakePathWithType(
		catgut.MakeFromString("project-2022-recurse"),
		catgut.MakeFromString("project-24q2-talent_show"),
	))

	e := catgut.MakeFromString("req-comp-internet")
	_, ok := es.All.ContainsEtikett(e)

	if !ok {
		t.Errorf("expected %s to be in %s", e, es)
	}
}

func BenchmarkMatchFirstYes(b *testing.B) {
	var es Etiketten

	es.AddPath(MakePathWithType(
		catgut.MakeFromString("area-home"),
		catgut.MakeFromString("project-reno"),
	))

	es.AddPath(MakePathWithType(
		catgut.MakeFromString("area-career"),
		catgut.MakeFromString("project-recurse"),
	))

	m := catgut.MakeFromString("area")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		es.All.ContainsEtikett(m)
	}
}

func BenchmarkMatchFirstNo(b *testing.B) {
	var es Etiketten

	es.AddPath(MakePathWithType(
		catgut.MakeFromString("area-home"),
		catgut.MakeFromString("project-reno"),
	))

	es.AddPath(MakePathWithType(
		catgut.MakeFromString("area-career"),
		catgut.MakeFromString("project-recurse"),
	))

	m := catgut.MakeFromString("x")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		es.All.ContainsEtikett(m)
	}
}
