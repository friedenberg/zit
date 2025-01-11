package repo_blobs

type Blob interface {
	GetRepoBlob() Blob
}
