package stringx

import (
	"unsafe"
)

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// toString performs unholy acts to avoid allocations
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
