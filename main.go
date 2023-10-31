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

	// generate design doc from folder
	designDoc, err := folderToJSON(parameters.source)
	if err != nil {
		fmt.Println("Error parsing folder:", err)
		os.Exit(1)
	}
	docId, idExists := designDoc["_id"]

	if !idExists {
		fmt.Println("document _id needs to be present")
		os.Exit(1)
	}

	docIdStr, ok := docId.(string)
	if !ok {
		fmt.Println("document _id needs to be present")
		os.Exit(1)
	}

	// should read revision of document, if available
	revision, err := getDocRevision(parameters.host, parameters.db, docIdStr, parameters.base64auth)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Revision:%s\n", revision)
	if revision != "" {
		// check to see if documents are identical
		identical, err := checkIdenticalDocs(parameters.host, parameters.db, designDoc, parameters.base64auth)
		if err != nil {
			fmt.Println("Error checking if docs are identical", err)
			os.Exit(1)
		}
		if identical {
			fmt.Println("Identical docs")
			os.Exit(1)
		}
		// otherwise set the existing revision needed for the document update
		designDoc["_rev"] = revision
	}

	// Convert the designDoc to a JSON string
	jsonString, err := json.MarshalIndent(designDoc, "", "    ")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonString))

	// push document
	uploadStatus, err := postDoc(parameters.host, parameters.db, jsonString, parameters.base64auth)
	if err != nil {
		fmt.Println("Error uploading document:", err)
		os.Exit(1)
	}

	fmt.Printf("Result status: %v", uploadStatus)
}
