package filehandling

import (
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

func GetAllFilesRecursiveInDirectory(filter func(string) bool, path string) ([]string, error) {
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
			// If it's a directory, make a recursive call
			subFiles, err := GetAllFilesRecursiveInDirectory(filter, fullPath)
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
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
