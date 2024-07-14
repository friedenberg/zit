package sku

import "fmt"

func (a *Transacted) String() string {
	return fmt.Sprintf(
		"%s %s %s",
		&a.ObjectId,
		a.GetObjectSha(),
		a.GetBlobSha(),
	)
}

func (a *Transacted) StringKennungBezeichnung() string {
	return fmt.Sprintf(
		"[%s %q]",
		&a.ObjectId,
		a.Metadata.Description,
	)
}

func (a *Transacted) StringKennungTai() string {
	return fmt.Sprintf(
		"%s@%s",
		&a.ObjectId,
		a.GetTai().StringDefaultFormat(),
	)
}

func (a *Transacted) StringKennungTaiAkte() string {
	return fmt.Sprintf(
		"%s@%s@%s",
		&a.ObjectId,
		a.GetTai().StringDefaultFormat(),
		a.GetBlobSha(),
	)
}

func (a *Transacted) StringKennungSha() string {
	return fmt.Sprintf(
		"%s@%s",
		&a.ObjectId,
		a.GetMetadata().Sha(),
	)
}

func (a *Transacted) StringKennungMutter() string {
	return fmt.Sprintf(
		"%s^@%s",
		&a.ObjectId,
		a.GetMetadata().Mutter(),
	)
}
