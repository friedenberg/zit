package umwelt

import (
	"bufio"
	"encoding/gob"
	"flag"
	"io"
	"os"
	"path"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/mutable_config"
	"code.linenisgreat.com/zit/go/zit/src/kilo/typ"
)

type Einleitung struct {
	AgeIdentity age.Identity
	Yin         string
	Yang        string
	Angeboren   immutable_config.Config
}

func (e *Einleitung) AddToFlagSet(f *flag.FlagSet) {
	f.Var(
		&e.AgeIdentity,
		"age",
		"",
	) // TODO-P3 move to Angeboren
	f.StringVar(&e.Yin, "yin", "", "File containing list of Kennung")
	f.StringVar(&e.Yang, "yang", "", "File containing list of Kennung")
	e.Angeboren.AddToFlagSet(f)
}

func (u *Umwelt) Einleitung(e Einleitung) (err error) {
	s := u.Standort()

	mkdirAll(s.DirKennung())
	mkdirAll(s.DirVerzeichnisse())
	mkdirAll(s.DirVerlorenUndGefunden())

	if err = readAndTransferLines(e.Yin, s.DirZit("Kennung", "Yin")); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = readAndTransferLines(e.Yang, s.DirZit("Kennung", "Yang")); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, g := range []genres.Genre{genres.Blob, genres.InventoryList} {
		var d string

		if d, err = s.DirObjektenGattung(
			e.Angeboren.GetStoreVersion(),
			g,
		); err != nil {
			if genres.IsErrUnsupportedGattung(err) {
				err = nil
				continue
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		mkdirAll(d)
	}

	if err = s.Age().AddIdentityOrGenerateIfNecessary(
		e.AgeIdentity,
		s.FileAge(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	writeFile(s.FileKonfigAngeboren(), e.Angeboren)

	if e.Angeboren.UseKonfigErworbenFile {
		writeFile(s.FileKonfigErworben(), "")
	} else {
		writeFile(s.FileKonfigCompiled(), "")
	}

	writeFile(s.FileSchlummernd(), "")

	if err = u.schlummernd.Flush(
		u.Standort(),
		u.PrinterHeader(),
		u.konfig.DryRun,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Standort().ResetVerzeichnisse(); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.TodoP2("determine if this should be an Einleitung option")
	if err = initDefaultTypAndKonfig(u); err != nil {
		err = errors.Wrap(err)
		return
	}

	{
		if err = u.Lock(); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, u.Unlock)

		if err = u.GetStore().ResetIndexes(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func initDefaultTypAndKonfig(u *Umwelt) (err error) {
	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	defaultTyp, defaultTypKennung := typ.Default()

	// var defaultTypTransacted *typ.Transacted

	var k ids.ObjectId

	if err = k.SetWithIdLike(defaultTypKennung); err != nil {
		err = errors.Wrap(err)
		return
	}

	{
		err = nil

		var sh interfaces.Sha

		if sh, _, err = u.GetStore().GetAkten().GetTypeV0().SaveBlobText(
			&defaultTyp,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		var k ids.ObjectId

		if err = k.SetWithIdLike(defaultTypKennung); err != nil {
			err = errors.Wrap(err)
			return
		}

		if _, err = u.GetStore().CreateOrUpdateAkteSha(
			&k,
			sh,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	{
		var sh interfaces.Sha

		if sh, err = writeDefaultErworben(u, defaultTypKennung); err != nil {
			err = errors.Wrap(err)
			return
		}

		if _, err = u.GetStore().CreateOrUpdateAkteSha(
			&ids.Config{},
			sh,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func writeDefaultErworben(
	u *Umwelt,
	dt ids.Type,
) (sh interfaces.Sha, err error) {
	defaultKonfig := mutable_config.Default(dt)

	f := u.GetStore().GetKonfigAkteFormat()

	var aw sha.WriteCloser

	if aw, err = u.Standort().BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	if _, err = f.FormatParsedBlob(aw, &defaultKonfig); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = sha.Make(aw.GetShaLike())

	return
}

func mkdirAll(elements ...string) {
	err := os.MkdirAll(path.Join(elements...), os.ModeDir|0o755)
	errors.PanicIfError(err)
}

func writeFile(p string, contents any) {
	var f *os.File
	var err error

	if f, err = files.CreateExclusiveWriteOnly(p); err != nil {
		if errors.IsExist(err) {
			ui.Err().Printf("%s already exists, not overwriting", p)
			err = nil
		} else {
		}

		return
	}

	defer errors.PanicIfError(err)
	defer errors.DeferredCloser(&err, f)

	if s, ok := contents.(string); ok {
		_, err = io.WriteString(f, s)
	} else {
		enc := gob.NewEncoder(f)
		err = enc.Encode(contents)
	}
}

func readAndTransferLines(in, out string) (err error) {
	errors.TodoP4("move to user operations")

	if in == "" {
		return
	}

	var fi, fo *os.File

	if fi, err = files.Open(in); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, fi.Close)

	if fo, err = files.CreateExclusiveWriteOnly(out); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, fo.Close)

	r := bufio.NewReader(fi)
	w := bufio.NewWriter(fo)

	defer errors.Deferred(&err, w.Flush)

	for {
		var l string
		l, err = r.ReadString('\n')

		if errors.Is(err, io.EOF) {
			err = nil
			break
		}

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		// TODO-P2 sterilize line
		w.WriteString(l)
	}

	return
}
