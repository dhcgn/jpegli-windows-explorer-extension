package install

import (
	"archive/zip"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhcgn/jpegli-windows-explorer-extension/types"
)

// Embedded zip files
var (
	//go:embed files/*
	files embed.FS
)

func Do() error {

	// Delete all folders in the application folder
	deleteAllFolders()

	// Copy Executable to the Program Data directory
	execPath := CopyExecutableToProgramData()
	if execPath == "" {
		fmt.Println("Failed to copy executable to Program Data directory")
		return fmt.Errorf("failed to copy executable to Program Data directory")
	}

	// Extract embedded zip files to the application folder
	_, err := ExtractEmbeddedZipFilesToAppFolder()
	if err != nil {
		fmt.Printf("Error extracting embedded zip files: %v\n", err)
		return fmt.Errorf("error extracting embedded zip files: %w", err)
	}

	// Set the executable as Windows Explorer context menu
	SetExecutableAsWindowsExplorerContextMenu(execPath)

	return nil
}

func deleteAllFolders() {
	// Get the application folder
	appFolder := getAppFolder()
	if appFolder == "" {
		fmt.Println("Failed to get application folder")
		return
	}

	// Delete all folders in the application folder
	err := os.RemoveAll(appFolder)
	if err != nil {
		fmt.Printf("Error deleting folders in %s: %v\n", appFolder, err)
	} else {
		fmt.Printf("Successfully deleted all folders in %s\n", appFolder)
	}
}

func getAppFolder() string {
	// Get the user cache directory
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		fmt.Printf("Error getting user cache directory: %v\n", err)
		return ""
	}
	// Create the target directory if it doesn't exist
	targetDir := filepath.Join(cacheDir, "jpegli-windows-explorer-extension")
	return targetDir
}

// GetToolsPath returns the paths to the exiftool and cjpegli executables
func GetToolsPath() (types.ExecutablePaths, error) {
	// Initialize empty paths
	execPaths := types.ExecutablePaths{
		Exiftool: "",
		Cjpegli:  "",
	}

	// Get the application folder
	appFolder := getAppFolder()
	if appFolder == "" {
		fmt.Println("Failed to get application folder")
		return execPaths, fmt.Errorf("failed to get application folder")
	}

	// Find the executables in the application folder
	execPaths.Exiftool = findExifTool(appFolder)
	execPaths.Cjpegli = findCjpegli(appFolder)

	if execPaths.Exiftool == "" {
		fmt.Println("Error: exiftool executable not found")
		return execPaths, fmt.Errorf("exiftool executable not found")
	}
	if execPaths.Cjpegli == "" {
		fmt.Println("Error: cjpegli executable not found")
		return execPaths, fmt.Errorf("cjpegli executable not found")
	}

	return execPaths, nil
}

// ExtractEmbeddedZipFilesToAppFolder extracts embedded zip files to the application folder
// Returns paths to the extracted executables
func ExtractEmbeddedZipFilesToAppFolder() (types.ExecutablePaths, error) {
	// Initialize empty paths
	execPaths := types.ExecutablePaths{
		Exiftool: "",
		Cjpegli:  "",
	}

	// Get the application folder
	appFolder := getAppFolder()
	if appFolder == "" {
		fmt.Println("Failed to get application folder")
		return execPaths, fmt.Errorf("failed to get application folder")
	}

	// Ensure the app folder exists
	err := os.MkdirAll(appFolder, 0755)
	if err != nil {
		fmt.Printf("Error creating application folder %s: %v\n", appFolder, err)
		return execPaths, fmt.Errorf("error creating application folder: %w", err)
	}

	// Read the embedded files directory
	entries, err := files.ReadDir("files")
	if err != nil {
		fmt.Printf("Error reading embedded files: %v\n", err)
		return execPaths, fmt.Errorf("error reading embedded files: %w", err)
	}

	// Extract each zip file
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".zip" {
			// Get the file path
			filePath := filepath.Join("files", entry.Name())

			fmt.Printf("Extracting '%s' ...\n", filePath)

			// Read the zip file content
			content, err := files.ReadFile("files/" + entry.Name())
			if err != nil {
				fmt.Printf("Error reading embedded file %s: %v\n", filePath, err)
				continue
			} // Create the target file
			targetPath := filepath.Join(appFolder, entry.Name())
			err = os.WriteFile(targetPath, content, 0644)
			if err != nil {
				fmt.Printf("Error writing file %s: %v\n", targetPath, err)
				continue
			}

			// Extract the zip file contents
			extractPath := filepath.Join(appFolder, filepath.Base(entry.Name()[:len(entry.Name())-4])) // Remove .zip extension
			err = os.MkdirAll(extractPath, 0755)
			if err != nil {
				fmt.Printf("Error creating extraction directory %s: %v\n", extractPath, err)
				continue
			} // Extract the zip file using Go's zip package
			fmt.Printf("Extracting contents from %s to %s\n", targetPath, extractPath)
			err = extractZipFile(targetPath, extractPath)
			if err != nil {
				fmt.Printf("Error extracting zip file %s: %v\n", targetPath, err)
				continue
			}

			fmt.Printf("Successfully extracted %s to %s\n", entry.Name(), targetPath)

			// Find and set executable paths based on zip file name
			if strings.Contains(entry.Name(), "exiftool") {
				// Find file 'exiftool(-k).exe' and rename it to 'exiftool.exe'
				exiftoolPath := findExifTool(extractPath)
				if exiftoolPath != "" {
					execPaths.Exiftool = exiftoolPath
				}
			} else if strings.Contains(entry.Name(), "jpegli") {
				// Find cjpegli.exe
				cjpegliPath := findCjpegli(extractPath)
				if cjpegliPath != "" {
					execPaths.Cjpegli = cjpegliPath
				}
			}
		}
	}

	if execPaths.Exiftool == "" {
		fmt.Println("Error: exiftool executable not found")
		return execPaths, fmt.Errorf("exiftool executable not found")
	}
	if execPaths.Cjpegli == "" {
		fmt.Println("Error: cjpegli executable not found")
		return execPaths, fmt.Errorf("cjpegli executable not found")
	}

	return execPaths, nil
}

// extractZipFile extracts the contents of a zip file to the specified destination directory
func extractZipFile(zipPath, destPath string) error {
	// Open the zip file
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer zipReader.Close()

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Extract each file
	for _, file := range zipReader.File {
		// Create full destination path
		dest := filepath.Join(destPath, file.Name)

		// Check for ZipSlip vulnerability
		if !strings.HasPrefix(dest, filepath.Clean(destPath)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", file.Name)
		}

		// If file is a directory, create it
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(dest, file.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dest, err)
			}
			continue
		}

		// Create directory for file if needed
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return fmt.Errorf("failed to create directory for file %s: %w", dest, err)
		}

		// Open destination file
		dstFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("failed to create destination file %s: %w", dest, err)
		}

		// Open source file from the zip
		srcFile, err := file.Open()
		if err != nil {
			dstFile.Close()
			return fmt.Errorf("failed to open file from zip %s: %w", file.Name, err)
		}

		// Copy content from zip to disk
		_, err = io.Copy(dstFile, srcFile)
		srcFile.Close()
		dstFile.Close()
		if err != nil {
			return fmt.Errorf("failed to copy file content from zip %s: %w", file.Name, err)
		}
	}

	return nil
}

// CopyExecutableToProgramData moves the executable to the Program Data directory
// Folder: os.UserCacheDir() + jpegli-windows-explorer-extension
// Return the path to the executable in the Program Data directory
func CopyExecutableToProgramData() string {
	// Get the current executable path
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting current executable path: %v\n", err)
		return ""
	}

	// Create the target directory if it doesn't exist
	targetDir := getAppFolder()
	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		fmt.Printf("Error creating directory %s: %v\n", targetDir, err)
		return ""
	}
	// Get the filename from the executable path
	targetPath := filepath.Join(targetDir, "jpegli-windows-explorer-extension.exe")

	// Always copy the executable to the target directory
	sourceFile, err := os.Open(exePath)
	if err != nil {
		fmt.Printf("Error opening source file %s: %v\n", exePath, err)
		return ""
	}
	defer sourceFile.Close()

	// Create the target file
	targetFile, err := os.Create(targetPath)
	if err != nil {
		fmt.Printf("Error creating target file %s: %v\n", targetPath, err)
		return ""
	}
	defer targetFile.Close()

	// Copy the contents
	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		fmt.Printf("Error copying file: %v\n", err)
		return ""
	}

	fmt.Printf("Executable successfully copied to %s\n", targetPath)
	return targetPath
}

// findExifTool searches for the exiftool executable in the given directory
// and renames exiftool(-k).exe to exiftool.exe if needed
func findExifTool(dir string) string {
	// Common patterns for exiftool executable
	patterns := []string{
		"exiftool.exe",
		"exiftool(-k).exe",
		"*/exiftool.exe",
		"*/exiftool(-k).exe",
		"**/exiftool.exe",
		"**/exiftool(-k).exe",
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			continue
		}

		for _, match := range matches {
			// If it's the exiftool(-k).exe variant, rename it to exiftool.exe
			if strings.Contains(match, "(-k)") {
				newPath := strings.Replace(match, "(-k)", "", 1)
				err := os.Rename(match, newPath)
				if err == nil {
					fmt.Printf("Renamed %s to %s\n", match, newPath)
					return newPath
				}
				fmt.Printf("Error renaming %s to %s: %v\n", match, newPath, err)
				return match // Return original path if rename fails
			}
			return match
		}
	}

	// Deeper recursive search if not found with glob
	var exiftoolPath string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && (strings.Contains(info.Name(), "exiftool") && strings.HasSuffix(info.Name(), ".exe")) {
			exiftoolPath = path
			return filepath.SkipAll
		}
		return nil
	})

	return exiftoolPath
}

// findCjpegli searches for the cjpegli executable in the given directory
func findCjpegli(dir string) string {
	// Common patterns for cjpegli executable
	patterns := []string{
		"cjpegli.exe",
		"*/cjpegli.exe",
		"bin/cjpegli.exe",
		"**/cjpegli.exe",
		"**/bin/cjpegli.exe",
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			continue
		}

		if len(matches) > 0 {
			return matches[0]
		}
	}

	// Deeper recursive search if not found with glob
	var cjpegliPath string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.Contains(info.Name(), "cjpegli") && strings.HasSuffix(info.Name(), ".exe") {
			cjpegliPath = path
			return filepath.SkipAll
		}
		return nil
	})

	return cjpegliPath
}
