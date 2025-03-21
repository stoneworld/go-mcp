package pkg

import (
	"log"
)

func Recover() {
	if r := recover(); r != nil {
		log.Println("Recovered in f", r)
	}
}
