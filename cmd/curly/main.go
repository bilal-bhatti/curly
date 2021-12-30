/*
Copyright © 2021 Bilal Bhatti
*/

package main

import (
	"log"

	"github.com/bilal-bhatti/curly/internal/curly"
)

var Version = "v0.0.0"

func main() {
	curly.Version = Version
	log.Println("curly", curly.Version)

	Execute()
}
