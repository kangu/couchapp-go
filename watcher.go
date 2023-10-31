package main

import (
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
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
