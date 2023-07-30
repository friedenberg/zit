package umwelt

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/mike/store_fs"
)

//		                                                _
//	   ___ ___  _ __ ___  _ __   ___  _ __   ___ _ __ | |_ ___
//	  / __/ _ \| '_ ` _ \| '_ \ / _ \| '_ \ / _ \ '_ \| __/ __|
//		| (_| (_) | | | | | | |_) | (_) | | | |  __/ | | | |_\__ \
//		 \___\___/|_| |_| |_| .__/ \___/|_| |_|\___|_| |_|\__|___/
//		                    |_|
func (u *Umwelt) FormatColorOptions() (o format.ColorOptions) {
	o.OffEntirely = !u.outIsTty
	return
}

func (u *Umwelt) FormatColorWriter() format.FuncColorWriter {
	if u.outIsTty {
		return format.MakeFormatWriterWithColor
	} else {
		return format.MakeFormatWriterNoopColor
	}
}

func (u *Umwelt) FormatSha(
	a func(sha.Sha) (string, error),
) schnittstellen.FuncWriterFormat[schnittstellen.ShaLike] {
	return kennung.MakeShaCliFormat(u.FormatColorWriter(), a)
}

//    ___  _     _      _    _
//   / _ \| |__ (_) ___| | _| |_ ___ _ __
//  | | | | '_ \| |/ _ \ |/ / __/ _ \ '_ \
//  | |_| | |_) | |  __/   <| ||  __/ | | |
//   \___/|_.__// |\___|_|\_\\__\___|_| |_|
//            |__/

func (u *Umwelt) FormatExternalFD() schnittstellen.FuncWriterFormat[kennung.FD] {
	return zettel.MakeCliFormatFD(
		u.Standort(),
		u.FormatColorWriter(),
	)
}

func (u *Umwelt) FormatFileNotRecognized() schnittstellen.FuncWriterFormat[kennung.FD] {
	return store_fs.MakeCliFormatNotRecognized(
		u.FormatColorWriter(),
		u.Standort(),
		u.FormatSha(u.StoreObjekten().GetAbbrStore().Shas().Abbreviate),
	)
}

func (u *Umwelt) FormatFDDeleted() schnittstellen.FuncWriterFormat[kennung.FD] {
	return store_fs.MakeCliFormatFDDeleted(
		u.Konfig().DryRun,
		u.FormatColorWriter(),
		u.Standort(),
		u.FormatExternalFD(),
	)
}
