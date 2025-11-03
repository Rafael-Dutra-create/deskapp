package main

import (
	"deskapp/src/internal/scripts"
	"log"
)

func main() {
    if err := scripts.MapTableToStruct(); err != nil {
        log.Fatal(err)
    }
}