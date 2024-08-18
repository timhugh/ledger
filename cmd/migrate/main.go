package main

import (
	"github.com/timhugh/ledger/db/sqlite"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing database file")
	}

	file := os.Args[1]
	log.Printf("opening database %s\n", file)
	client, err := sqlite.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("running migrations\n")

	err = client.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("migrations complete\n")
}
