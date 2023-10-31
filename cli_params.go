package main

import (
	"errors"
	"flag"
)

type cliParams struct {
	host   string
	db     string
	source string
	watch  bool
}

func readCLIParams() (cliParams, error) {
	// Get the non-flag arguments (i.e., DB in this case)
	host := flag.String("host", "http://localhost:5984", "CouchDB host address")
	db := flag.String("db", "", "Target database")
	source := flag.String("source", ".", "Source directory")
	watch := flag.Bool("watch", false, "Live folder watch")
	flag.Parse()

	if *db == "" {
		return cliParams{}, errors.New("target db needs to be specified with -db option")
	}

	return cliParams{
		host:   *host,
		db:     *db,
		source: *source,
		watch:  *watch,
	}, nil
}
