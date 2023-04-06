package objekte

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

/*
transacted format <- transacted
external objekte format <- sku.ExternalMaybe, Metadatei, T
objekte format <- Metadatei, T
akte format <- T
*/

type AkteFormat[A any, APtr schnittstellen.Ptr[A]] interface {
	ParseAkte(io.Reader, APtr) (int64, schnittstellen.Sha, error)
	FormatAkte(io.Writer, APtr) (int64, schnittstellen.Sha, error)
}

type ObjekteFormat[A any, APtr schnittstellen.Ptr[A]] interface {
	ParseObjekte(io.Reader, metadatei.Metadatei, APtr) (int64, schnittstellen.Sha, error)
	FormatObjekte(io.Writer, metadatei.Metadatei, APtr) (int64, schnittstellen.Sha, error)
}
