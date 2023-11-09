package konfig

import (
	"encoding/gob"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/delta/typ_akte"
)

func (kc *Compiled) loadKonfigErworben(s standort.Standort) (err error) {
	var f *os.File

	p := s.FileKonfigCompiled()

	if kc.UseKonfigErworbenFile {
		p = s.FileKonfigErworben()
	}

	if f, err = files.Open(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	dec := gob.NewDecoder(f)

	if err = dec.Decode(&kc.compiled); err != nil {
		if errors.IsEOF(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	return
}

func (kc *Compiled) Flush(
	s standort.Standort,
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) (err error) {
	if !kc.hasChanges || kc.DryRun {
		return
	}

	if err = kc.recompile(tagp); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := s.FileKonfigCompiled()

	if kc.UseKonfigErworbenFile {
		p = s.FileKonfigErworben()
	}

	var f *os.File

	if f, err = files.OpenExclusiveWriteOnlyTruncate(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	dec := gob.NewEncoder(f)

	if err = dec.Encode(kc.compiled); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
