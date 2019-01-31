package trick

import (
	"reflect"
	"unsafe"
)

// Equal compares string and bytes
func Equal(a string, b []byte) bool {
	return a == *(*string)(unsafe.Pointer(&b))
}

// String2Bytes converts string to bytes
func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// Bytes2String converts bytes to string
func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
