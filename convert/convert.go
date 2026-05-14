package convert

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dhcgn/jpegli-windows-explorer-extension/types"
)

type ConvertStats struct {
	FileSizeRatio float64
	TargetSize    int64
	SourceSize    int64
	SavedSize     int64
}

const OptimizedByTag = "XMP-jpegli:OptimizedBy"

func Convert(tools types.ExecutablePaths, distance float64, overrideOriginal bool, sourcePath, targetPath, markerValue string) (ConvertStats, error) {
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
	var tempFile *os.File
	if overrideOriginal {
		// Use a temporary file if we're going to override the original
		// Create temp file in the same directory as source to ensure we're on the same filesystem
		dir := filepath.Dir(sourcePath)
		var err error
		tempFile, err = os.CreateTemp(dir, ".jpegli-*.tmp")
		if err != nil {
			return ConvertStats{}, fmt.Errorf("failed to create temporary file: %w", err)
		}
		tempFile.Close() // Close immediately, we just need the path
		actualTargetPath = tempFile.Name()
	}

	// Use exec.Command to run cjpegli with the provided distance parameter
	cmd := exec.Command(tools.Cjpegli, sourcePath, actualTargetPath, "-d", fmt.Sprintf("%.1f", distanceValue))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ConvertStats{}, fmt.Errorf("cjpegli execution failed: %w\nOutput: %s", err, output)
	}

	// Step 2: Copy metadata from source to target using ExifTool
	copyMetadataArgs := withExiftoolConfig(tools, "-overwrite_original", "-TagsFromFile", sourcePath, actualTargetPath)
	cmd = exec.Command(tools.Exiftool, copyMetadataArgs...)
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

	// Step 4: Mark target file as optimized only after conversion and metadata copy succeeded.
	markerErr := MarkAsOptimized(tools, actualTargetPath, markerValue)
	if markerErr != nil {
		return ConvertStats{}, markerErr
	}

	// Step 5: Get file sizes and calculate statistics
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

func ReadOptimizedBy(tools types.ExecutablePaths, sourcePath string) (string, error) {
	if tools.Exiftool == "" {
		return "", fmt.Errorf("exiftool path is empty")
	}
	if tools.ExiftoolConfig == "" {
		return "", fmt.Errorf("exiftool config path is empty")
	}

	// -q -q for quiet mode to suppress warnings
	args := withExiftoolConfig(tools, "-s3", "-q", "-q", "-"+OptimizedByTag, sourcePath)
	cmd := exec.Command(tools.Exiftool, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("exiftool read marker failed: %w\nOutput: %s", err, output)
	}
	return strings.TrimSpace(string(output)), nil
}

func MarkAsOptimized(tools types.ExecutablePaths, targetPath, markerValue string) error {
	if tools.Exiftool == "" {
		return fmt.Errorf("exiftool path is empty")
	}
	if tools.ExiftoolConfig == "" {
		return fmt.Errorf("exiftool config path is empty")
	}

	tagAssignment := fmt.Sprintf("-%s=%s", OptimizedByTag, markerValue)
	args := withExiftoolConfig(tools, "-overwrite_original", tagAssignment, targetPath)
	cmd := exec.Command(tools.Exiftool, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("exiftool write marker failed: %w\nOutput: %s", err, output)
	}
	return nil
}

func withExiftoolConfig(tools types.ExecutablePaths, args ...string) []string {
	capacity := len(args)
	if tools.ExiftoolConfig != "" {
		capacity += 2
	}
	withConfig := make([]string, 0, capacity)
	if tools.ExiftoolConfig != "" {
		withConfig = append(withConfig, "-config", tools.ExiftoolConfig)
	}
	withConfig = append(withConfig, args...)
	return withConfig
}
