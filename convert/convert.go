package convert

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dhcgn/jpegli-windows-explorer-extension/types"
)

type ConvertStats struct {
	FileSizeRatio float64
	TargetSize    int64
	SourceSize    int64
	SavedSize     int64
}

func Convert(tools types.ExecutablePaths, distance float64, overrideOriginal bool, sourcePath, targetPath string) (ConvertStats, error) {
	// Validate tools paths
	if tools.Cjpegli == "" {
		return ConvertStats{}, fmt.Errorf("cjpegli path is empty")
	}
	if tools.Exiftool == "" {
		return ConvertStats{}, fmt.Errorf("exiftool path is empty")
	}

	// Check if the source file exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return ConvertStats{}, fmt.Errorf("source file doesn't exist: %s", sourcePath)
	}

	// Get the source file size before conversion
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return ConvertStats{}, fmt.Errorf("error getting source file info: %w", err)
	}
	sourceSize := sourceInfo.Size()
	// Step 1: Convert image with cjpegli.exe using the distance option
	// Default to 1.0 if not specified (visually lossless)
	// Allowed range is 0.0 to 25.0
	distanceValue := 1.0
	if distance >= 0.0 && distance <= 25.0 {
		distanceValue = distance
	} else if distance > 25.0 {
		distanceValue = 25.0
	} else if distance < 0.0 {
		distanceValue = 0.0
	}

	// Determine the actual target path (temporary or final)
	actualTargetPath := targetPath
	if overrideOriginal {
		// Use temporary file if we're going to override the original
		actualTargetPath = sourcePath + ".jpegli.tmp"
	}

	// Use exec.Command to run cjpegli with the provided distance parameter
	cmd := exec.Command(tools.Cjpegli, sourcePath, actualTargetPath, "-d", fmt.Sprintf("%.1f", distanceValue))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ConvertStats{}, fmt.Errorf("cjpegli execution failed: %w\nOutput: %s", err, output)
	}

	// Step 2: Copy metadata from source to target using ExifTool
	cmd = exec.Command(tools.Exiftool, "-overwrite_original", "-TagsFromFile", sourcePath, actualTargetPath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		// Clean up temporary file if it exists
		if overrideOriginal {
			os.Remove(actualTargetPath)
		}
		return ConvertStats{}, fmt.Errorf("exiftool execution failed: %w\nOutput: %s", err, output)
	}

	// Step 3: If overrideOriginal is true and both tools succeeded, replace the original file
	if overrideOriginal {
		// Both cjpegli and exiftool have succeeded, now replace the original
		err = os.Rename(actualTargetPath, sourcePath)
		if err != nil {
			os.Remove(actualTargetPath) // Clean up temporary file
			return ConvertStats{}, fmt.Errorf("failed to replace original file: %w", err)
		}
		// Update targetPath to sourcePath for stats calculation
		actualTargetPath = sourcePath
	}

	// Step 4: Get file sizes and calculate statistics
	targetInfo, err := os.Stat(actualTargetPath)
	if err != nil {
		return ConvertStats{}, fmt.Errorf("error getting target file info: %w", err)
	}
	targetSize := targetInfo.Size()

	// Calculate ratio (target size / source size)
	ratio := float64(targetSize) / float64(sourceSize)

	return ConvertStats{
		FileSizeRatio: ratio,
		SourceSize:    sourceSize,
		TargetSize:    targetSize,
		SavedSize:     sourceSize - targetSize,
	}, nil
}
