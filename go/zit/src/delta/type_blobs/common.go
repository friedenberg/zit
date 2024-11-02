package type_blobs

type Common interface {
	GetFileExtension() string
	GetBinary() bool
}
