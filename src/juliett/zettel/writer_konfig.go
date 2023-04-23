package zettel

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/india/konfig"
)

func MakeWriterKonfig(
	k konfig.Compiled,
) schnittstellen.FuncIter[*Transacted] {
	errors.TodoP0("switch to sigils")
	errors.TodoP3("add efficient parsing of hiding tags")

	if k.IncludeHidden {
		return collections.MakeWriterNoop[*Transacted]()
	}

	return func(z *Transacted) (err error) {
		for _, p := range z.Verzeichnisse.Etiketten.Sorted {
			for _, t := range k.EtikettenHiddenStringsSlice {
				if strings.HasPrefix(p, t) {
					err = collections.MakeErrStopIteration()
					return
				}
			}
		}

		t := k.GetApproximatedTyp(z.GetTyp()).ApproximatedOrActual()

		if t != nil && t.Objekte.Akte.Archived {
			err = collections.MakeErrStopIteration()
			return
		}

		return
	}
}
