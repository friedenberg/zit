package objekte

import (
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/matcher"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type ExternalLike interface {
	GetFDs() sku.ExternalFDs
	GetObjekteFD() kennung.FD
	GetAkteFD() kennung.FD
	matcher.Matchable
	metadatei.Getter
	GetSkuLike() sku.SkuLike
}

type ExternalLikePtr interface {
	ExternalLike
	GetFDsPtr() *sku.ExternalFDs
	metadatei.GetterPtr
	metadatei.Setter
	GetKennungLikePtr() kennung.KennungPtr
	GetSkuLikePtr() sku.SkuLikePtr
}
