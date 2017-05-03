package main

import (
	"github.com/trichner/evedt"
	"log"
)

func main() {
	//TODO should pass config here instead of magic in evedt package
	log.Fatal(evedt.Start())
}
