package main

import (
	"fmt"
	"os"
)

const version = "1.0.0"

func main() {
	parameters, err := ReadCLIParams()
	if err != nil {
		// fmt.Println(err)
		printUsageInstructions()
		os.Exit(1)
	}
	if parameters.showVersion {
		fmt.Printf("%s\n", version)
		os.Exit(1)
	}

	ProcessFolderToCouch(parameters)

	if parameters.watch {
		KeepWatchingForChanges(parameters)
	}
}

func printUsageInstructions() {
	fmt.Println("Usage: couchapp-go [options]")
	fmt.Println()
	fmt.Println("Description: Command-line tool for pushing a folder structure into a CouchDB design document.")
	fmt.Println()
	fmt.Println("Options")
	fmt.Println("--db=[dbname]		Target database")
	fmt.Println("--source=[path]		Path for couchapp folder. Defaults to current folder")
	fmt.Println("--host=[addr]		IP/address of CouchDB server. Defaults to http://localhost:5984")
	fmt.Println("--user=[name]		Username for basic authentication (if needed)")
	fmt.Println("--pass=[pass]		Password for basic authentication (if needed)")
	fmt.Println("--watch			If set to true, watch source folder and push on file changes")
}
