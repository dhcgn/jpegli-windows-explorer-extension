package filehandling

import (
	"fmt"
	"os"
)

func IsPathDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if info.IsDir() {
		return true, nil
	}
	return false, nil
}

func GetAllFilesInDirectory(filter func(string) bool, path string, warn func(string)) ([]string, error) {
	var allFiles []string

	isDir, err := IsPathDir(path)
	if err != nil {
		return nil, err
	}

	if !isDir {
		// It's a file, apply the filter the same way we do for files in directories
		if filter == nil || filter(path) {
			allFiles = append(allFiles, path)
		}
		return allFiles, nil
	}

	// Get all files in the directory recursively
	var files []string

	// Read directory contents
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	// Process each entry in the directory
	for _, entry := range entries {
		fullPath := path + string(os.PathSeparator) + entry.Name()

		if entry.IsDir() {
			warn(fmt.Sprintf("Skipping is directory: %s", fullPath))
		} else {
			// If it's a file, check if it passes the filter
			if filter == nil || filter(fullPath) {
				files = append(files, fullPath)
			}
		}
	}

	allFiles = append(allFiles, files...)
	return allFiles, nil
}
