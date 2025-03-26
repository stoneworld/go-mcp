package pkg

import (
	"context"
	"log"
	"runtime/debug"
)

func Recover() {
	if r := recover(); r != nil {
		log.Println("Recovered in f", r, string(debug.Stack()))
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
