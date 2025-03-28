package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		_, err := fmt.Fprintf(os.Stdout, "%s\n", scanner.Bytes())
		if err != nil {
			return
		}
	}
}
