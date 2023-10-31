package main

import (
	"os"
	"path/filepath"
	"strings"
)

// folderToJSON recursively converts a folder structure to a JSON object.
func folderToJSON(folderPath string) (map[string]interface{}, error) {
	folderJSON := make(map[string]interface{})

	// Read the contents of the folder
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filePath := filepath.Join(folderPath, file.Name())

		if file.IsDir() {
			// If it's a subfolder, recursively call folderToJSON
			subfolderJSON, err := folderToJSON(filePath)
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
			folderJSON[key] = string(content)
		}
	}

	return folderJSON, nil
}
