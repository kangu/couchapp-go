package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FolderToJSON recursively converts a folder structure to a JSON object.
func FolderToJSON(folderPath string) (map[string]interface{}, error) {
	folderJSON := make(map[string]interface{})

	// Read the contents of the folder
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filePath := filepath.Join(folderPath, file.Name())

		// if filename starts with ., just ignore it
		if strings.HasPrefix(file.Name(), ".") {
			// fmt.Printf("Ignoring file: %v\n", file.Name())
			continue
		}

		if file.IsDir() {
			// If it's a subfolder, recursively call folderToJSON
			subfolderJSON, err := FolderToJSON(filePath)
			if err != nil {
				return nil, err
			}
			folderJSON[file.Name()] = subfolderJSON
		} else {
			// If it's a file, read its contents and add to the JSON
			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, err
			}
			// Remove the file extension from the key
			key := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			// Skip files with no name, just extension
			// like .couchapprc or .DS_Store
			if key != "" {
				folderJSON[key] = string(content)
			}
		}
	}

	return folderJSON, nil
}

func ProcessFolderToCouch(parameters CliParams) {
	// generate design doc from folder
	designDoc, err := FolderToJSON(parameters.source)
	if err != nil {
		fmt.Println("Error parsing folder:", err)
		os.Exit(1)
	}
	docId, idExists := designDoc["_id"]

	if !idExists {
		folderPath, err := filepath.Abs(parameters.source)
		if err != nil {
			fmt.Println("Error resolving absolute path:", err)
			return
		}
		fmt.Printf("Couchapp not detected at %s, file named \"_id\" needs to be present\n", folderPath)
		fmt.Printf("Use --source parameter to specify a different path\n")
		os.Exit(1)
	}

	docIdStr, ok := docId.(string)
	if !ok {
		// should never happen, by design what's read from disk is always a string
		os.Exit(1)
	}

	// should read revision of document, if available
	revision, err := getDocRevision(parameters.host, parameters.db, docIdStr, parameters.base64auth)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if revision != "" {
		// check to see if documents are identical
		identical, err := checkIdenticalDocs(parameters.host, parameters.db, designDoc, parameters.base64auth)
		if err != nil {
			fmt.Println("Error checking if docs are identical", err)
			os.Exit(1)
		}
		if identical {
			fmt.Printf("Identical docs, rev %s not changed\n", revision)
			return
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
	// fmt.Printf("Uploading %v", string(jsonString))

	// push document
	uploadStatus, err := postDoc(parameters.host, parameters.db, jsonString, parameters.base64auth)
	if err != nil {
		fmt.Println("Error uploading document:", err)
		return
	}

	if uploadStatus == 201 {
		// get revision again
		revision, err = getDocRevision(parameters.host, parameters.db, docIdStr, parameters.base64auth)
		fmt.Printf("Successful upload, new rev %s\n", revision)
	} else {
		fmt.Printf("Upload status error: %v\n", uploadStatus)
	}

}
