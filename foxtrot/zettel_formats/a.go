package zettel_formats

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/akte_ext"
	"github.com/friedenberg/zit/bravo/files"
	"github.com/friedenberg/zit/bravo/line_format"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/delta/objekte"
	"github.com/friedenberg/zit/echo/zettel"
)

type (
	_Sha     = sha.Sha
	_Etikett = etikett.Etikett

	_Zettel                   = zettel.Zettel
	_ZettelFormat             = zettel.Format
	_ZettelFormatContextRead  = zettel.FormatContextRead
	_ZettelFormatContextWrite = zettel.FormatContextWrite
	_AkteWriterFactory        = zettel.AkteWriterFactory

	_EtikettSet    = etikett.Set
	_AkteExt       = akte_ext.AkteExt
	_ObjekteWriter = objekte.Writer
)

var (
	_EtikettNewSet       = etikett.NewSet
	_LineFormatNewWriter = line_format.NewWriter
	_Errorf              = errors.Errorf
	_PanicIfError        = errors.PanicIfError
	_Error               = errors.Error
	_Open                = open_file_guard.Open
	_Create              = open_file_guard.Create
	_Close               = open_file_guard.Close
	_FilesExists         = files.Exists
)
