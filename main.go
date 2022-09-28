package main

import (
	"flag"
	"time"
)

func GetCurrentTimestamp() int64 {
	return time.Now().UTC().Unix()
}

func main() {
	id := flag.String("id", "", "identity to use")
	test := flag.String("t", "", "specifies which test to run")

	flag.Parse()

	if *test != "" {
		if *test == "transaction" {
			TestTransaction()
		}
	} else {
		if *id != "" {
			id := LoadIdentity(*id)
			GenesisBootstrap(id)
		} else {
			panic("no id specified")
		}
	}

}
