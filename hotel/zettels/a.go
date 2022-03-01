package zettels

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/node_type"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/files"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/age"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/charlie/konfig"
	"github.com/friedenberg/zit/delta/objekte"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/echo/sharded_store"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/akten"
	"github.com/friedenberg/zit/foxtrot/etiketten"
	"github.com/friedenberg/zit/foxtrot/hinweisen"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
	"github.com/friedenberg/zit/golf/stored_zettel_formats"
)

const (
	_TypeAkte = node_type.TypeAkte
)

type (
	_Age = age.Age

	_Etikett              = etikett.Etikett
	_EtikettExpanderRight = etikett.ExpanderRight
	_Etiketten            = etiketten.Etiketten

	_Hinweis   = hinweis.Hinweis
	_Hinweisen = hinweisen.Hinweisen

	_Id = id.Id

	_Sha = sha.Sha

	_Entry        = sharded_store.Entry
	_Shard        = sharded_store.Shard
	_ShardGeneric = sharded_store.ShardGeneric
	_Store        = sharded_store.Store

	_Zettel                   = zettel.Zettel
	_ZettelFormat             = zettel.Format
	_ZettelFormatContextRead  = zettel.FormatContextRead
	_ZettelFormatContextWrite = zettel.FormatContextWrite
	_AkteReaderFactory        = zettel.AkteReaderFactory
	_AkteWriterFactory        = zettel.AkteWriterFactory

	_StoredZettel              = stored_zettel.Stored
	_NamedZettel               = stored_zettel.Named
	_ZettelExternal            = stored_zettel.External
	_ZettelCheckedOut          = stored_zettel.CheckedOut
	_NamedZettelFilter         = stored_zettel.NamedFilter
	_StoredZettelFormatObjekte = stored_zettel_formats.Objekte

	_Akten              = akten.Akten
	_ErrorDuplicateAtke = akten.DuplicateAkteError

	_ObjekteReader = objekte.Reader
	_ObjekteWriter = objekte.Writer

	_ZettelFormatText = zettel_formats.Text

	_Konfig = konfig.Konfig

	_Umwelt = umwelt.Umwelt
)

var (
	_Close                  = open_file_guard.Close
	_Create                 = open_file_guard.Create
	_Errf                   = stdprinter.Errf
	_Outf                   = stdprinter.Outf
	_Error                  = errors.Error
	_ErrorAs                = errors.As
	_FilesExist             = files.Exists
	_ErrorsIs               = errors.Is
	_Errorf                 = errors.Errorf
	_ErrorNormal            = errors.Normal
	_EtikettNewSet          = etikett.NewSet
	_IdHeadTailFromFileName = id.HeadTailFromFileName
	_IdMakeDirNecessary     = id.MakeDirIfNecessary
	_IdPath                 = id.Path
	_MakeBlindHinweis       = hinweis.MakeBlindHinweis
	_NewAkten               = akten.New
	_NewEtiketten           = etiketten.New
	_NewHinweisen           = hinweisen.New
	_NewShard               = sharded_store.NewShard
	_NewStore               = sharded_store.NewStore
	_ObjekteDecodeBase64    = objekte.DecodeBase64
	_ObjekteEncodeBase64    = objekte.EncodeBase64
	_ObjekteNewReader       = objekte.NewReader
	_ObjekteNewWriterMover  = objekte.NewWriterMover
	_ObjekteRead            = objekte.Read
	_ObjekteWriteAndMove    = objekte.WriteAndMove
	_Open                   = open_file_guard.Open
	_OpenFile               = open_file_guard.OpenFile
	_PanicIfError           = stdprinter.PanicIfError
	_ReadDirNames           = open_file_guard.ReadDirNames
	_ShaFromHash            = sha.FromHash
)
