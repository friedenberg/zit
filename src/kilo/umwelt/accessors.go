package umwelt

import (
	"io"

	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/echo/akten"
	"github.com/friedenberg/zit/src/india/store_objekten"
	"github.com/friedenberg/zit/src/juliett/store_working_directory"
	"github.com/friedenberg/zit/src/juliett/zettel_printer"
	"github.com/friedenberg/zit/src/zettel_verzeichnisse"
)

func (u *Umwelt) Konfig() konfig.Konfig {
	return u.konfig
}

func (u *Umwelt) PrinterOut() *zettel_printer.Printer {
	return u.printerOut
}

func (u *Umwelt) Out() io.Writer {
	return u.out
}

func (u *Umwelt) In() io.Reader {
	return u.in
}

func (u *Umwelt) Err() io.Writer {
	return u.err
}

func (u *Umwelt) Standort() standort.Standort {
	return u.standort
}

func (u *Umwelt) StoreAkten() akten.Akten {
	return u.akten
}

func (u *Umwelt) StoreObjekten() *store_objekten.Store {
	return u.storeObjekten
}

func (u *Umwelt) StoreWorkingDirectory() *store_working_directory.Store {
	return u.storeWorkingDirectory
}

func (u *Umwelt) ZettelVerzeichnissePool() *zettel_verzeichnisse.Pool {
	return u.zettelVerzeichnissePool
}
