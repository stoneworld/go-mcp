package pkg

import (
	"log"
	"runtime/debug"
)

func Recover() {
	if r := recover(); r != nil {
		log.Println("Recovered in f", r, string(debug.Stack()))
	}
}
