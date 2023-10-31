package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func watchDirRecursively(watcher *fsnotify.Watcher, dirPath string) error {
	err := watcher.Add(dirPath)
	if err != nil {
		return err
	}

	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			err := watcher.Add(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func KeepWatchingForChanges(parameters CliParams) {
	folderPath, err := filepath.Abs(parameters.source)
	if err != nil {
		fmt.Println("Error resolving absolute path:", err)
		return
	}
	fmt.Println("Watching folder", folderPath)
	for {
		// Create a new watcher
		// reset on every change so that it picks up new subfolders being created
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

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
						ProcessFolderToCouch(parameters)
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
