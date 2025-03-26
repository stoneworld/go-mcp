package pkg

import (
	"log"
	"runtime/debug"
	"unsafe"
)

func Recover() {
	if r := recover(); r != nil {
		log.Printf("panic: %v\nstack: %s", r, debug.Stack())
	}
}

func RecoverWithFunc(f func(r any)) {
	if r := recover(); r != nil {
		f(r)
		log.Printf("panic: %v\nstack: %s", r, debug.Stack())
	}
}

func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
