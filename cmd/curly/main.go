/*
Copyright © 2021 Bilal Bhatti
*/

package main

import "log"

var Version = "v0.0.0-DEV"

func main() {
	log.Println("version:", Version)

	Execute()
}
