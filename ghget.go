// Package main
package main

import (
	"fmt"
	"log"

	"github.com/ta7eralla/ghget/flags"
)

const (
	publicURL  = "https://raw.githubusercontent.com/%s/%s/refs/heads/%s/%s"
	privateURL = ""
)

func main() {
	log.SetPrefix("ghget: ")
	log.SetFlags(0)

	fcfg, err := flags.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(fcfg)
}
