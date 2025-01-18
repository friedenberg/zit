package interfaces

import (
	"flag"
	"io"
)

type ReadWrapper interface {
	WrapReader(r io.Reader) (io.ReadCloser, error)
}

type WriteWrapper interface {
	WrapWriter(w io.Writer) (io.WriteCloser, error)
}

type IOWrapper interface {
	ReadWrapper
	WriteWrapper
}

type BlobCompression interface {
	flag.Value
	IOWrapper
	GetBlobCompression() BlobCompression
}

type BlobEncryption interface {
	flag.Value
	IOWrapper
	GetBlobEncryption() BlobEncryption
}

type BlobReader interface {
	BlobReader(Sha) (ShaReadCloser, error)
}

type BlobWriter interface {
	BlobWriter() (ShaWriteCloser, error)
}

type BlobStore interface {
	GetBlobStore() BlobStore
	HasBlob(sh Sha) (ok bool)
	BlobReader
	BlobWriter
}

type BlobStoreConfigImmutable interface {
	GetBlobStoreConfigImmutable() BlobStoreConfigImmutable
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
)
