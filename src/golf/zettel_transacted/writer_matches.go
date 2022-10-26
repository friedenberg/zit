package zettel_transacted

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/golf/zettel_external"
)

type WriterMatchesOptions struct {
	ShasZettelen sha.Set
	ShasAkten    sha.Set
}

func MakeWriterMatchesOptions(zes ...zettel_external.Zettel) WriterMatchesOptions {
	sz := sha.MakeMutableSet()
	sa := sha.MakeMutableSet()

	for _, z := range zes {
		sz.Add(z.Named.Stored.Sha)
		akteSha := z.Named.Stored.Zettel.Akte

		if !akteSha.IsNull() {
			sa.Add(z.Named.Stored.Zettel.Akte)
		}
	}

	return WriterMatchesOptions{
		ShasZettelen: sz.Copy(),
		ShasAkten:    sa.Copy(),
	}
}

type WriterMatchesReasons map[sha.Sha]hinweis.Set

func (wmr *WriterMatchesReasons) Add(sh sha.Sha, h hinweis.Hinweis) {
	var hsm hinweis.MutableSet
	hs, ok := (*wmr)[sh]

	if !ok {
		hsm = hinweis.MakeMutableSet()
	} else {
		hsm = hs.MutableCopy()
	}

	hsm.Add(h)

	(*wmr)[sh] = hsm.Copy()
}

type WriterMatches struct {
	options              WriterMatchesOptions
	matchReasonsZettelen WriterMatchesReasons
	matchReasonsAkten    WriterMatchesReasons
	found                Set
}

func MakeWriterMatches(options WriterMatchesOptions) WriterMatches {
	return WriterMatches{
		options:              options,
		matchReasonsZettelen: make(map[sha.Sha]hinweis.Set),
		matchReasonsAkten:    make(map[sha.Sha]hinweis.Set),
		found:                MakeSetHinweis(0),
	}
}

func (w *WriterMatches) WriteZettelTransacted(z *Zettel) (err error) {
	switch {
	case w.found.Contains(*z):
		old, ok := w.found.GetString(w.found.GetKey(*z))

		if !ok || old.Schwanz.Less(z.Schwanz) {
			w.found.Add(*z)
		}

	case w.options.ShasZettelen.Contains(z.Named.Stored.Sha):
		w.found.Add(*z)
		w.matchReasonsAkten.Add(z.Named.Stored.Sha, z.Named.Hinweis)

	case w.options.ShasAkten.Contains(z.Named.Stored.Zettel.Akte):
		w.found.Add(*z)
		w.matchReasonsAkten.Add(z.Named.Stored.Zettel.Akte, z.Named.Hinweis)

	default:
		err = io.EOF
		return
	}

	return
}

func (w WriterMatches) MatchReasonsZettelen() WriterMatchesReasons {
	return w.matchReasonsZettelen
}

func (w WriterMatches) MatchReasonsAkten() WriterMatchesReasons {
	return w.matchReasonsAkten
}

func (w WriterMatches) Matches() Set {
	return w.found
}
