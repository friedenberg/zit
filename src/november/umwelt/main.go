package umwelt

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/file_lock"
	"github.com/friedenberg/zit/src/charlie/kennung"
	"github.com/friedenberg/zit/src/delta/konfig"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	"github.com/friedenberg/zit/src/juliett/zettel_verzeichnisse"
	"github.com/friedenberg/zit/src/lima/store_objekten"
	"github.com/friedenberg/zit/src/mike/store_fs"
)

type Umwelt struct {
	in  *os.File
	out *os.File
	err *os.File

	inIsTty  bool
	outIsTty bool
	errIsTty bool

	standort standort.Standort
	konfig   konfig.Konfig

	storesInitialized     bool
	lock                  *file_lock.Lock
	storeObjekten         *store_objekten.Store
	age                   *age.Age
	storeWorkingDirectory *store_fs.Store

	zettelVerzeichnissePool zettel_verzeichnisse.Pool
}

func Make(kCli konfig.Cli) (u *Umwelt, err error) {
	u = &Umwelt{
		in:                      os.Stdin,
		out:                     os.Stdout,
		err:                     os.Stderr,
		zettelVerzeichnissePool: zettel_verzeichnisse.MakePool(),
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
	return u.Initialize(u.Konfig().Cli)
}

func (u *Umwelt) Initialize(kCli konfig.Cli) (err error) {
	if err = u.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

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

	if u.konfig, err = konfig.Make(u.standort.FileKonfigToml(), kCli); err != nil {
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
		u.age = &age.Age{}
		// if u.age, err = age.MakeDefaultTest(); err != nil {
		// 	errors.Wrap(err)
		// 	return
		// }
	}

	for _, rb := range u.konfig.Recipients {
		if err = u.age.AddBech32PivYubikeyEC256(rb); err != nil {
			errors.Wrap(err)
			return
		}
	}

	u.storeObjekten, err = store_objekten.Make(
		u.lock,
		*u.age,
		u.konfig,
		u.standort,
		u.zettelVerzeichnissePool,
	)

	if err != nil {
		err = errors.Wrapf(err, "failed to initialize zettel meta store")
		return
	}

	csk := store_fs.Konfig{
		Konfig:       u.konfig,
		CacheEnabled: u.konfig.CheckoutCacheEnabled,
	}

	errors.Print("initing checkout store")
	u.storeWorkingDirectory, err = store_fs.New(
		csk,
		u.standort,
		u.storeObjekten,
	)

	if err != nil {
		errors.Print(err)
		err = errors.Wrap(err)
		return
	}

	errors.Print("done initing checkout store")

	u.storeObjekten.SetZettelTransactedLogWriter(
		store_objekten.ZettelTransactedLogWriters{
			New:       u.PrinterZettelTransacted(format.StringNew),
			Updated:   u.PrinterZettelTransacted(format.StringUpdated),
			Unchanged: u.PrinterZettelTransacted(format.StringUnchanged),
			Archived:  u.PrinterZettelTransacted(format.StringArchived),
		},
	)

	u.storeWorkingDirectory.SetZettelCheckedOutWriters(
		store_fs.ZettelCheckedOutLogWriters{
			ZettelOnly: u.PrinterZettelCheckedOutFresh(zettel_checked_out.ModeZettelOnly),
			AkteOnly:   u.PrinterZettelCheckedOutFresh(zettel_checked_out.ModeZettelOnly),
			Both:       u.PrinterZettelCheckedOutFresh(zettel_checked_out.ModeZettelOnly),
		},
	)

	u.storesInitialized = true

	return
}

func (u Umwelt) DefaultEtiketten() (etiketten kennung.EtikettSet, err error) {
	metiketten := kennung.MakeEtikettMutableSet()

	for e, t := range u.konfig.Tags {
		if !t.AddToNewZettels {
			continue
		}

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
