package interfaces

import (
	"flag"
	"io"
)

type BlobIOMiddleware interface {
	WrapReader(r io.Reader) (io.ReadCloser, error)
	WrapWriter(w io.Writer) (io.WriteCloser, error)
}

type BlobCompression interface {
	flag.Value
	BlobIOMiddleware
  GetBlobCompression() BlobCompression
}

type BlobEncryption interface {
	flag.Value
	BlobIOMiddleware
  GetBlobEncryption() BlobEncryption
}

type BlobStore interface {
	GetBlobStore() BlobStore
	HasBlob(sh Sha) (ok bool)
	BlobWriter() (w ShaWriteCloser, err error)
	BlobReader(sh Sha) (r ShaReadCloser, err error)
}

type BlobStoreConfig interface {
	GetBlobStoreImmutableConfig() BlobStoreConfig
	GetBlobEncryption() BlobEncryption
	GetBlobCompression() BlobCompression
	GetLockInternalFiles() bool
}

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
		BlobReader(Sha) (ShaReadCloser, error)
	}

	BlobWriterFactory interface {
		BlobWriter() (ShaWriteCloser, error)
	}
)
