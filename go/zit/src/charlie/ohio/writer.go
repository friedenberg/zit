package ohio

import (
	"encoding/binary"
	"io"
	"unsafe"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

func WriteAllOrDieTrying(w io.Writer, b []byte) (n int, err error) {
	var acc int

	for n < len(b) {
		acc, err = w.Write(b[n:])
		n += acc
		if err != nil {
			return
		}
	}

	return
}

func WriteInt8(w io.Writer, n int8) (written int, err error) {
	b := [1]byte{byte(n)}

	written, err = WriteAllOrDieTrying(w, b[:])
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func WriteUint8(w io.Writer, n uint8) (written int, err error) {
	b := [1]byte{n}

	written, err = WriteAllOrDieTrying(w, b[:])
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func WriteUint16(w io.Writer, n uint16) (written int, err error) {
	var intErr int
	var b [binary.MaxVarintLen16]byte

	intErr = binary.PutUvarint(b[:], uint64(n))

	if intErr != 1 {
		err = errors.Errorf("expected to write %d but wrote %d", 2, intErr)
		return
	}

	written, err = WriteAllOrDieTrying(w, b[:intErr])
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func WriteUint32(w io.Writer, n uint32) (written int, err error) {
	var intErr int
	var b [binary.MaxVarintLen32]byte

	intErr = binary.PutVarint(b[:], int64(n))

	if intErr != 1 {
		err = errors.Errorf("expected to write %d but wrote %d", 2, intErr)
		return
	}

	written, err = WriteAllOrDieTrying(w, b[:intErr])
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func WriteInt64(w io.Writer, n int64) (written int, err error) {
	var intErr int
	var b [binary.MaxVarintLen64]byte

	intErr = binary.PutVarint(b[:], n)

	written, err = WriteAllOrDieTrying(w, b[:intErr])
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func WriteFixedInt64(w io.Writer, n int64) (written int, err error) {
	b := Int64ToByteArray(n)

	written, err = WriteAllOrDieTrying(w, b[:])
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func Int64ToByteArray(i int64) [8]byte {
	return *(*[unsafe.Sizeof(i)]byte)(unsafe.Pointer(&i))
}

func ByteArrayToInt64(arr [8]byte) int64 {
	val := int64(0)
	size := len(arr)

	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}

	return val
}
