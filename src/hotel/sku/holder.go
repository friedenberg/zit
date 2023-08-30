package sku

import (
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type Holder struct {
	Metadatei   metadatei.Metadatei
	KennungLike kennung.Kennung
}

func (h *Holder) GetMetadatei() metadatei.Metadatei {
	return h.Metadatei
}

func (h *Holder) SetMetadatei(m metadatei.Metadatei) {
	h.Metadatei = m
}

func (h *Holder) SetKennungLike(kl kennung.Kennung) (err error) {
	h.KennungLike = kl
	return
}
