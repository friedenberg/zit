package umwelt

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/file_lock"
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/echo/akten"
	"github.com/friedenberg/zit/src/hotel/store_objekten"
	"github.com/friedenberg/zit/src/india/store_working_directory"
	"github.com/friedenberg/zit/src/juliett/zettel_printer"
)

type Umwelt struct {
	standort standort.Standort
	konfig   konfig.Konfig

	logger errors.Logger

	in  *os.File
	out *os.File
	err *os.File

	storesInitialized     bool
	lock                  *file_lock.Lock
	storeObjekten         *store_objekten.Store
	akten                 akten.Akten
	age                   age.Age
	storeWorkingDirectory *store_working_directory.Store
	printerOut            *zettel_printer.Printer
}

func Make(c konfig.Konfig) (u *Umwelt, err error) {
	u = &Umwelt{
		konfig: c,
		logger: c.Logger,
		in:     os.Stdin,
		out:    os.Stdout,
		err:    os.Stderr,
	}

	err = u.Initialize()

	return
}

func (u *Umwelt) Initialize() (err error) {
	if u.standort, err = standort.Make(u.konfig); err != nil {
		err = errors.Wrap(err)
		return
	}

	u.lock = file_lock.New(u.standort.DirZit("Lock"))

	fa := u.standort.FileAge()

	if files.Exists(fa) {
		if u.age, err = age.Make(fa); err != nil {
			errors.Wrap(err)
			return
		}
	} else {
		u.age = age.MakeEmpty()
	}

	if u.storeObjekten, err = store_objekten.Make(u.lock, u.age, u.konfig, u.standort); err != nil {
		err = errors.Wrapf(err, "failed to initialize zettel meta store")
		return
	}

	if u.akten, err = akten.New(u.standort.DirZit()); err != nil {
		err = errors.Wrap(err)
		return
	}

	csk := store_working_directory.Konfig{
		Konfig:       u.konfig,
		CacheEnabled: u.konfig.CheckoutCacheEnabled,
	}

	errors.Print("initing checkout store")

	if u.storeWorkingDirectory, err = store_working_directory.New(csk, u.standort.Cwd(), u.storeObjekten); err != nil {
		errors.Print(err)
		err = errors.Wrap(err)
		return
	}

	errors.Print("done initing checkout store")

	u.printerOut = zettel_printer.Make(u.konfig, u.storeObjekten, u.out)
	//TODO move to konfig
	// u.printerOut.ShouldAbbreviateHinweisen = true

	u.storesInitialized = true

	return
}

func (u Umwelt) DefaultEtiketten() (etiketten etikett.Set, err error) {
	etiketten = etikett.MakeSet()

	for e, t := range u.konfig.Tags {
		if !t.AddToNewZettels {
			continue
		}

		if err = etiketten.AddString(e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
