package interfaces

// Blobs represent persisted files, like blobs in Git. Blobs are used by
// Zettels, types, tags, config, and inventory lists.
type (
	Blob[T any] interface{}

	BlobPtr[T any] interface {
		Blob[T]
		Ptr[T]
	}

	BlobGetter[
		V any,
	] interface {
		GetBlob(Sha) (V, error)
	}

	BlobPutter[
		V any,
	] interface {
		PutBlob(V)
	}

	BlobGetterPutter[
		V any,
	] interface {
		BlobGetter[V]
		BlobPutter[V]
	}

	BlobIOFactory interface {
		BlobReaderFactory
		BlobWriterFactory
	}

	BlobReaderFactory interface {
		BlobReader(ShaGetter) (ShaReadCloser, error)
	}

	BlobWriterFactory interface {
		BlobWriter() (ShaWriteCloser, error)
	}
)
