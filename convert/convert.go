package convert

import (
	"io"
	"os"

	"github.com/dhcgn/jpegli-windows-explorer-extention/types"
)

type ConvertOptions struct {
}

type ConvertStats struct {
	FileSizeRatio float64
	TargetSize    int64
	SourceSize    int64
	SavedSize     int64
}

func Convert(tools types.ExecutablePaths, opts ConvertOptions, sourcePath, targetPath string) (ConvertStats, error) {
	// copy file from sourcePath to targetPath
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return ConvertStats{}, err
	}
	defer sourceFile.Close()

	// Create the target file
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return ConvertStats{}, err
	}
	defer targetFile.Close()

	// Copy the contents from source to target
	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		return ConvertStats{}, err
	}

	// Get file sizes to calculate ratio
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return ConvertStats{}, err
	}

	targetInfo, err := os.Stat(targetPath)
	if err != nil {
		return ConvertStats{}, err
	}
	// Calculate ratio (target size / source size)
	ratio := float64(targetInfo.Size()) / float64(sourceInfo.Size())

	return ConvertStats{
		FileSizeRatio: ratio,
		SourceSize:    sourceInfo.Size(),
		TargetSize:    targetInfo.Size(),
		SavedSize:     sourceInfo.Size() - targetInfo.Size(),
	}, nil
}
