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

func (a *Transacted) StringObjectIdDescription() string {
	return fmt.Sprintf(
		"[%s %q]",
		&a.ObjectId,
		a.Metadata.Description,
	)
}

func (a *Transacted) StringObjectIdTai() string {
	return fmt.Sprintf(
		"%s@%s",
		&a.ObjectId,
		a.GetTai().StringDefaultFormat(),
	)
}

func (a *Transacted) StringObjectIdTaiBlob() string {
	return fmt.Sprintf(
		"%s@%s@%s",
		&a.ObjectId,
		a.GetTai().StringDefaultFormat(),
		a.GetBlobSha(),
	)
}

func (a *Transacted) StringObjectIdSha() string {
	return fmt.Sprintf(
		"%s@%s",
		&a.ObjectId,
		a.GetMetadata().Sha(),
	)
}

func (a *Transacted) StringObjectIdParent() string {
	return fmt.Sprintf(
		"%s^@%s",
		&a.ObjectId,
		a.GetMetadata().Mutter(),
	)
}
