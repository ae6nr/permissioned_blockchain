package main

import (
	"errors"
	"fmt"
	"os"
)

func printUsage() {
	fmt.Println("Permissioned Blockchain")
	fmt.Println("Please specify a command.")
	fmt.Println("  serve")
	fmt.Println("     starts a blockchain server")
	fmt.Println("  entry <server_url> <identity> <entry>")
	fmt.Println("     submit an entry transaction with the data <entry> to the server <server_url> using your <identity>")
}

// check that there are at least n command line arguments after the program name
// if not, print usage
func checkOsArgs(n int) error {
	if len(os.Args) < n+1 {
		printUsage()
		return errors.New("not enough arguments")
	}
	return nil
}

func main() {

	if err := checkOsArgs(1); err != nil {
		return
	}

	cmd := os.Args[1]

	if cmd == "serve" {
		GetCurrentTimestamp()
		StartServer()
	} else if cmd == "entry" {
		if err := checkOsArgs(4); err != nil {
			return
		}
		server_url := os.Args[2]
		id := LoadIdentity(os.Args[3])
		entry := []byte(os.Args[4])

		SubmitEntry(server_url, entry, id)
	} else if cmd == "bootstrap" {
		if err := checkOsArgs(2); err != nil {
			return
		}
		fmt.Printf("Creating a genesis block for %s\r\n", os.Args[2])
		GenesisBootstrap(LoadIdentity(os.Args[2]))
	} else {
		printUsage()
	}
}
