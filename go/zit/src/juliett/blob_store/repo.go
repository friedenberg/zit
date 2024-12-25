package blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_blobs"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_layout"
)

type RepoStore struct {
	dirLayout repo_layout.Layout
	v0        Store[repo_blobs.V0, *repo_blobs.V0]
}

func MakeRepoStore(
	dirLayout repo_layout.Layout,
) RepoStore {
	return RepoStore{
		dirLayout: dirLayout,
		v0: MakeBlobStore(
			dirLayout,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[repo_blobs.V0](
					dirLayout,
				),
				ParsedBlobTomlFormatter[repo_blobs.V0, *repo_blobs.V0]{},
				dirLayout,
			),
			func(a *repo_blobs.V0) {
				a.Reset()
			},
		),
	}
}
