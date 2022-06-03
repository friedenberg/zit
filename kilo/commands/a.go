package commands

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/node_type"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/akte_ext"
	"github.com/friedenberg/zit/bravo/files"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/age"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/charlie/file_lock"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/charlie/konfig"
	"github.com/friedenberg/zit/charlie/script_value"
	"github.com/friedenberg/zit/delta/objekte"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/etiketten"
	"github.com/friedenberg/zit/foxtrot/hinweisen"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
	"github.com/friedenberg/zit/golf/alfred"
	"github.com/friedenberg/zit/golf/organize_text"
	"github.com/friedenberg/zit/golf/stored_zettel_formats"
	"github.com/friedenberg/zit/hotel/zettels"
)

const (
	_TypeAkte    = node_type.TypeAkte
	_TypeEtikett = node_type.TypeEtikett
	_TypeHinweis = node_type.TypeHinweis
	_TypeUnknown = node_type.TypeUnknown
	_TypeZettel  = node_type.TypeZettel
)

type (
	_Id = id.Id

	_ObjekteWriter = objekte.Writer

	_OrganizeChanges       = organize_text.Changes
	_OrganizeText          = organize_text.Text
	_OrganizeTextErrorRead = organize_text.ErrorRead
	_OrganizeTextOptions   = organize_text.Options

	_Age = age.Age

	_AkteExt = akte_ext.AkteExt

	_Hinweis   = hinweis.Hinweis
	_Hinweisen = hinweisen.Hinweisen
	_Provider  = hinweis.Provider

	_ErrorsStackTracer = errors.StackTracer

	_Sha = sha.Sha

	_Type = node_type.Type

	_Umwelt = *umwelt.Umwelt

	_Etiketten            = etiketten.Etiketten
	_Etikett              = etikett.Etikett
	_EtikettExpanderRight = etikett.ExpanderRight
	_EtikettSet           = etikett.Set

	_ScriptValue = script_value.ScriptValue

	_AlfredWriter = alfred.Writer

	_ZettelCheckedOut         = stored_zettel.CheckedOut
	_ExternalZettel           = stored_zettel.External
	_ZettelsCheckinOptions    = zettels.CheckinOptions
	_Zettels                  = zettels.Zettels
	_ZettelsChain             = zettels.Chain
	_VerlorenAndGefundenError = zettels.VerlorenAndGefundenError

	_RemoteScript     = konfig.RemoteScript
	_RemoteScriptFile = konfig.RemoteScriptFile
	_Konfig           = konfig.Konfig

	_ZettelFormatContextWrite = zettel.FormatContextWrite
	_Zettel                   = zettel.Zettel
	_ZettelFormat             = zettel.Format
	_ZettelFormatContextRead  = zettel.FormatContextRead

	_StoredZettel               = stored_zettel.Stored
	_NamedZettel                = stored_zettel.Named
	_StoredZettelFormatsObjekte = stored_zettel_formats.Objekte

	_ZettelFormatsText                        = zettel_formats.Text
	_ZettelFormatTextErrAkteInlineAndFilePath = zettel_formats.ErrHasInlineAkteAndFilePath

	_FileLock = file_lock.Lock
)

var (
	_AgeGenerate            = age.Generate
	_AlfredNewWriter        = alfred.NewWriter
	_Close                  = open_file_guard.Close
	_Create                 = open_file_guard.Create
	_DeleteFilesAndDirs     = open_file_guard.DeleteFilesAndDirs
	_Err                    = stdprinter.Err
	_Errf                   = stdprinter.Errf
	_ErrorAs                = errors.As
	_EtikettNewSet          = etikett.NewSet
	_FilesExist             = files.Exists
	_HinweisNewEmpty        = hinweis.NewEmpty
	_IdHeadTailFromFileName = id.HeadTailFromFileName
	_IdMakeDirNecessary     = id.MakeDirIfNecessary
	_IdPath                 = id.Path
	_KonfigDefaultCli       = konfig.DefaultCli
	_MakeBlindHinweis       = hinweis.MakeBlindHinweis
	_MakeBlindHinweisParts  = hinweis.MakeBlindHinweisParts
	_NewEtiketten           = etiketten.New
	_NewHinweisen           = hinweisen.New
	_NewZettels             = zettels.New
	_ObjekteRead            = objekte.Read
	_ObjekteWrite           = objekte.Write
	_ObjekteWriteAndMove    = objekte.WriteAndMove
	_Open                   = open_file_guard.Open
	_ReadDirNames           = open_file_guard.ReadDirNames
	_OpenEditor             = open_file_guard.OpenEditor
	_OpenVimWithArgs        = open_file_guard.OpenVimWithArgs
	_OpenFile               = open_file_guard.OpenFile
	_OpenFiles              = open_file_guard.OpenFiles
	_OrganizeTextNew        = organize_text.New
	_OrganizeTextNewEmpty   = organize_text.NewEmpty
	_Out                    = stdprinter.Out
	_Outf                   = stdprinter.Outf
	_PanicIfError           = stdprinter.PanicIfError
	_TempFile               = open_file_guard.TempFile
	_TempFileWithPattern    = open_file_guard.TempFileWithPattern
)
