package convert

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dhcgn/jpegli-windows-explorer-extension/types"
)

type ConvertOptions struct {
	Distance float64
}

type ConvertStats struct {
	FileSizeRatio float64
	TargetSize    int64
	SourceSize    int64
	SavedSize     int64
}

func Convert(tools types.ExecutablePaths, opts ConvertOptions, sourcePath, targetPath string) (ConvertStats, error) {
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
	distance := 1.0
	if opts.Distance >= 0.0 && opts.Distance <= 25.0 {
		distance = opts.Distance
	} else if opts.Distance > 25.0 {
		distance = 25.0
	} else if opts.Distance < 0.0 {
		distance = 0.0
	}

	// Use exec.Command to run cjpegli with the provided distance parameter
	cmd := exec.Command(tools.Cjpegli, sourcePath, targetPath, "-d", fmt.Sprintf("%.1f", distance))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ConvertStats{}, fmt.Errorf("cjpegli execution failed: %w\nOutput: %s", err, output)
	}

	// Step 2: Copy metadata from source to target using ExifTool
	cmd = exec.Command(tools.Exiftool, "-overwrite_original", "-TagsFromFile", sourcePath, targetPath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return ConvertStats{}, fmt.Errorf("exiftool execution failed: %w\nOutput: %s", err, output)
	}

	// Step 3: Get file sizes and calculate statistics
	targetInfo, err := os.Stat(targetPath)
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
