package main

import (
	"encoding/base64"
	"errors"
	"flag"
)

type cliParams struct {
	host       string
	db         string
	base64auth string
	source     string
	watch      bool
}

func readCLIParams() (cliParams, error) {
	// Get the non-flag arguments (i.e., DB in this case)
	host := flag.String("host", "http://localhost:5984", "CouchDB host address")
	db := flag.String("db", "", "Target database")
	user := flag.String("user", "", "Username")
	pass := flag.String("pass", "", "Password")
	source := flag.String("source", ".", "Source directory")
	watch := flag.Bool("watch", false, "Live folder watch")
	flag.Parse()

	if *db == "" {
		return cliParams{}, errors.New("target db needs to be specified with -db option")
	}

	authHeader := ""
	if (*user != "") && (*pass != "") {
		authHeader = "Basic " + base64.StdEncoding.EncodeToString([]byte(*user+":"+*pass))
	}

	return cliParams{
		host:       *host,
		db:         *db,
		base64auth: authHeader,
		source:     *source,
		watch:      *watch,
	}, nil
}
