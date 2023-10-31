package main

import (
	"fmt"
	"os"
)

func main() {
	parameters, err := readCLIParams()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	processFolderToCouch(parameters)

	if parameters.watch {
		keepWatchingForChanges(parameters)
	}
}
