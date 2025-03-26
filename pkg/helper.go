package pkg

import (
	"context"
	"log"
	"runtime/debug"
	"unsafe"
)

func Recover() {
	if r := recover(); r != nil {
		log.Printf("panic: %v\nstack: %s", r, debug.Stack())
	}
}

func SafeRunGo(ctx context.Context, f func()) {
	defer func(ctx context.Context) {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r, string(debug.Stack()))
		}
	}(ctx)

	f()
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
