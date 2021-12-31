/*
Copyright © 2021 Bilal Bhatti
*/

package main

import (
	"github.com/bilal-bhatti/curly/internal/curly"
)

var Version = "0.0.0"

func main() {
	curly.Version = Version
	Execute()
}
