package main

import (
	"log"

	"code.google.com/p/rsc/devweb/slave"
	_ "github.com/coreos/discovery/http"
)

func main() {
	log.SetFlags(0)
	slave.Main()
}
