package ohio

import "unsafe"

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

func Int32ToByteArray(i int32) [4]byte {
	return *(*[unsafe.Sizeof(i)]byte)(unsafe.Pointer(&i))
}

func ByteArrayToInt32(arr [4]byte) int32 {
	val := int32(0)
	size := len(arr)

	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}

	return val
}
