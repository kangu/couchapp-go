package main

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {

	parameters, err := readCLIParams()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	processFolderToCouch(parameters)

	if parameters.watch {
		for {
			// Create a new watcher
			// reset on every change so it picks up new subfolders created
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				log.Fatal(err)
			}
			defer watcher.Close()

			folderPath, err := filepath.Abs(parameters.source)
			if err != nil {
				fmt.Println("Error resolving absolute path:", err)
				return
			}
			fmt.Println("Watching folder", folderPath)
			// Watch the root folder and its subdirectories recursively
			err = watchDirRecursively(watcher, folderPath)
			if err != nil {
				log.Fatal(err)
			}

			// Create a channel to debounce events
			debounceDuration := 500 * time.Millisecond
			events := make(chan fsnotify.Event)
			go func() {
				var timer *time.Timer
				for {
					select {
					case event := <-events:
						if timer != nil {
							timer.Stop()
						}
						timer = time.AfterFunc(debounceDuration, func() {
							log.Println("Event:", event)
							processFolderToCouch(parameters)
						})
					}
				}
			}()

			// Start an infinite loop to wait for events
		eventsLoop:
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						break eventsLoop // Break out of the outer loop
					}

					// Check if the file name should be ignored (e.g., .DS_Store on OSX)
					if strings.Contains(event.Name, ".DS_Store") {
						continue
					}

					// Check if the event is a CHMOD event and ignore it
					if event.Op&fsnotify.Chmod == fsnotify.Chmod {
						continue
					}

					// Handle the file system event
					events <- event // Send the event to the debounce channel

					// Check for folder creation events
					if event.Op&fsnotify.Create == fsnotify.Create {
						break eventsLoop // Break out of the outer loop
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						break eventsLoop // Break out of the outer loop
					}
					fmt.Println("Error:", err)
				}
			}
		}
	}
}

func processFolderToCouch(parameters cliParams) {
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
	fmt.Printf("Uploading %v", string(jsonString))

	// push document
	uploadStatus, err := postDoc(parameters.host, parameters.db, jsonString, parameters.base64auth)
	if err != nil {
		fmt.Println("Error uploading document:", err)
		os.Exit(1)
	}

	if uploadStatus == 201 {
		// get revision again
		revision, err = getDocRevision(parameters.host, parameters.db, docIdStr, parameters.base64auth)
		fmt.Printf("Successful upload, new rev %s\n", revision)
	} else {
		fmt.Printf("Upload status status: %v\n", uploadStatus)
	}

}
