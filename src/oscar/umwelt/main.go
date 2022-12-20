package umwelt

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/file_lock"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/echo/standort"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/india/konfig"
	"github.com/friedenberg/zit/src/juliett/konfig_compiled"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/store_objekten"
	"github.com/friedenberg/zit/src/november/store_fs"
)

type Umwelt struct {
	sonnenaufgang ts.Time

	in  *os.File
	out *os.File
	err *os.File

	inIsTty  bool
	outIsTty bool
	errIsTty bool

	standort standort.Standort
	konfig   konfig_compiled.Compiled

	storesInitialized     bool
	lock                  *file_lock.Lock
	storeObjekten         *store_objekten.Store
	age                   *age.Age
	storeWorkingDirectory *store_fs.Store

	zettelVerzeichnissePool *zettel.PoolVerzeichnisse
}

func Make(kCli konfig.Cli) (u *Umwelt, err error) {
	u = &Umwelt{
		in:                      os.Stdin,
		out:                     os.Stdout,
		err:                     os.Stderr,
		zettelVerzeichnissePool: collections.MakePool[zettel.Verzeichnisse](),
	}

	if files.IsTty(u.in) {
		u.inIsTty = true
	}

	if files.IsTty(u.out) {
		u.outIsTty = true
	}

	if files.IsTty(u.err) {
		u.errIsTty = true
	}

	err = u.Initialize(kCli)

	return
}

func (u *Umwelt) Reset() (err error) {
	return u.Initialize(u.Konfig().Cli())
}

func (u *Umwelt) Initialize(kCli konfig.Cli) (err error) {
	if err = u.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	u.sonnenaufgang = ts.Now()

	//TODO-P4 consider moving to konfig_compiled
	{
		if kCli.Verbose {
			errors.SetVerbose()
		} else {
			log.SetOutput(ioutil.Discard)
		}

		standortOptions := standort.Options{
			BasePath: kCli.BasePath,
		}

		if standortOptions.BasePath == "" {
			if standortOptions.BasePath, err = os.Getwd(); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if u.standort, err = standort.Make(standortOptions); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	{
		fa := u.standort.FileAge()

		if files.Exists(fa) {
			if u.age, err = age.Make(fa); err != nil {
				errors.Wrap(err)
				return
			}
		} else {
			u.age = &age.Age{}
			// if u.age, err = age.MakeDefaultTest(); err != nil {
			// 	errors.Wrap(err)
			// 	return
			// }
		}
	}

	{
		var k *konfig_compiled.Compiled

		if k, err = konfig_compiled.Make(
			u.standort,
			kCli,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		u.konfig = *k
	}

	u.lock = file_lock.New(u.standort.DirZit("Lock"))

	// for _, rb := range u.konfig.Transacted.Objekte.Akte.Recipients {
	// 	if err = u.age.AddBech32PivYubikeyEC256(rb); err != nil {
	// 		errors.Wrap(err)
	// 		return
	// 	}
	// }

	if u.storeObjekten, err = store_objekten.Make(
		u.lock,
		*u.age,
		u.KonfigPtr(),
		u.standort,
		u.zettelVerzeichnissePool,
	); err != nil {
		err = errors.Wrapf(err, "failed to initialize zettel meta store")
		return
	}

	errors.Log().Print("initing checkout store")
	if u.storeWorkingDirectory, err = store_fs.New(
		u.Sonnenaufgang(),
		u.Konfig(),
		u.standort,
		u.storeObjekten,
	); err != nil {
		errors.Log().Print(err)
		err = errors.Wrap(err)
		return
	}

	errors.Log().Print("done initing checkout store")

	u.storeObjekten.Zettel().SetZettelTransactedLogWriter(
		u.ZettelTransactedLogPrinters(),
	)

	u.storeObjekten.Konfig().SetKonfigLogWriters(
		store_objekten.KonfigLogWriters{
			Updated:   u.PrinterKonfigTransacted(format.StringUpdated),
			Unchanged: u.PrinterKonfigTransacted(format.StringUnchanged),
		},
	)

	u.storeObjekten.Typ().SetTypLogWriters(
		store_objekten.TypLogWriters{
			New:       u.PrinterTypTransacted(format.StringUpdated),
			Updated:   u.PrinterTypTransacted(format.StringUpdated),
			Unchanged: u.PrinterTypTransacted(format.StringUnchanged),
		},
	)

	u.storeObjekten.Etikett().SetEtikettLogWriters(
		store_objekten.EtikettLogWriters{
			New:       u.PrinterEtikettTransacted(format.StringUpdated),
			Updated:   u.PrinterEtikettTransacted(format.StringUpdated),
			Unchanged: u.PrinterEtikettTransacted(format.StringUnchanged),
		},
	)

	u.storeWorkingDirectory.SetZettelExternalLogPrinter(
		u.PrinterZettelExternal(),
	)

	u.storesInitialized = true

	return
}

// TODO-P2 remove this
func (u Umwelt) DefaultEtiketten() (etiketten kennung.EtikettSet, err error) {
	metiketten := kennung.MakeEtikettMutableSet()

	for _, e := range u.konfig.EtikettenToAddToNew {
		if err = metiketten.AddString(e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	etiketten = metiketten.Copy()

	return
}

func (u Umwelt) Flush() error {
	return u.age.Close()
}
