package file_extensions

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
)

// TODO make new non-german version
type FileExtensions struct {
	Zettel   string `toml:"zettel"`
	Organize string `toml:"organize"`
	Type     string `toml:"typ"`
	Tag      string `toml:"etikett"`
	Repo     string `toml:"kasten"`
}

func (a FileExtensions) GetFileExtensionForGenre(
	g1 interfaces.GenreGetter,
) string {
	g := genres.Must(g1)

	switch g {
	case genres.Zettel:
		return a.GetFileExtensionZettel()

	case genres.Type:
		return a.GetFileExtensionType()

	case genres.Tag:
		return a.GetFileExtensionTag()

	case genres.Repo:
		return a.GetFileExtensionRepo()

	default:
		return ""
	}
}

func (a FileExtensions) GetFileExtensionGetter() interfaces.FileExtensionGetter {
	return a
}

func (a FileExtensions) GetFileExtensionZettel() string {
	return a.Zettel
}

func (a FileExtensions) GetFileExtensionOrganize() string {
	return a.Organize
}

func (a FileExtensions) GetFileExtensionType() string {
	return a.Type
}

func (a FileExtensions) GetFileExtensionTag() string {
	return a.Tag
}

func (a FileExtensions) GetFileExtensionRepo() string {
	return a.Repo
}

func (a *FileExtensions) Reset() {
	a.Zettel = ""
	a.Organize = ""
	a.Type = ""
	a.Tag = ""
	a.Repo = ""
}

func (a *FileExtensions) ResetWith(b FileExtensions) {
	a.Zettel = b.Zettel
	a.Organize = b.Organize
	a.Type = b.Type
	a.Tag = b.Tag
	a.Repo = b.Repo
}
