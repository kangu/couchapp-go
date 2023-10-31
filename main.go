package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {

	parameters, err := readCLIParams()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%v", parameters)

	folderPath := "./sample_couchapp" // Replace with the path to your folder
	result, err := folderToJSON(folderPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Convert the result to a JSON string
	jsonString, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonString))
}
