package main

import (
	"encoding/base64"
	"errors"
	"flag"
)

type CliParams struct {
	showVersion bool
	host        string
	db          string
	base64auth  string
	source      string
	watch       bool
}

func ReadCLIParams() (CliParams, error) {
	showVersion := flag.Bool("v", false, "Print the version number")
	host := flag.String("host", "http://localhost:5984", "CouchDB host address")
	db := flag.String("db", "", "Target database")
	user := flag.String("user", "", "Username")
	pass := flag.String("pass", "", "Password")
	source := flag.String("source", ".", "Source directory")
	watch := flag.Bool("watch", false, "Live folder watch")
	flag.Parse()

	if *showVersion {
		return CliParams{
			showVersion: true,
		}, nil
	}

	// at least DB parameter needs to be specified
	if *db == "" {
		return CliParams{}, errors.New("target db needs to be specified with --db option")
	}

	authHeader := ""
	if (*user != "") && (*pass != "") {
		authHeader = "Basic " + base64.StdEncoding.EncodeToString([]byte(*user+":"+*pass))
	}

	return CliParams{
		host:       *host,
		db:         *db,
		base64auth: authHeader,
		source:     *source,
		watch:      *watch,
	}, nil
}
