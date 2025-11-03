// cmd/createapp/main.go
package main

import (
	"deskapp/src/internal/scripts"
	"log"
)

func main() {
    if err := scripts.CreateApp(); err != nil {
        log.Fatal(err)
    }
}