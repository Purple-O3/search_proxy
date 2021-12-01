package tools

import (
	"reflect"
	"time"
	"unsafe"
)

func Str2Bytes(s string) []byte {
	return (*[0x7fff0000]byte)(unsafe.Pointer(
		(*reflect.StringHeader)(unsafe.Pointer(&s)).Data),
	)[:len(s):len(s)]
}

func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func TimeCost() func() time.Duration {
	start := time.Now()
	return func() time.Duration {
		return time.Since(start)
	}
}
