package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhcgn/jpegli-windows-explorer-extension/convert"
	"github.com/dhcgn/jpegli-windows-explorer-extension/filehandling"
	"github.com/dhcgn/jpegli-windows-explorer-extension/install"
	"github.com/pterm/pterm"
)

const (
	AppName = "jpegli-windows-explorer-extension"
)

var (
	Version = "UNSET"
	Build   = "UNSET"
	Commit  = "UNSET"
)

func waitForAnyKey() {
	fmt.Println("Press any key to continue...")
	var input [1]byte
	os.Stdin.Read(input[:])
}

func main() {
	fmt.Println("jpegli-windows-explorer-extension")
	fmt.Printf("Version: %s, Build: %s, Commit: %s\n", Version, Build, Commit)

	if len(os.Args) == 1 {
		pterm.Println("No arguments provided. Want to install and set context menu?")
		result, _ := pterm.DefaultInteractiveConfirm.Show()
		pterm.Println()
		pterm.Info.Printfln("You answered: %s", boolToText(result))
		if result {
			install.Do()
			pterm.Println("Installation completed.")
		} else {
			pterm.Println("You chose not to install.")
		}
		waitForAnyKey()
		return
	}

	// Print all args, excluded the first one (the program name)
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("file/folder: %d: %s\n", i, os.Args[i])
	}

	// For the moment we only support one file/folder at a time
	if len(os.Args) > 2 {
		pterm.Error.Printfln("Only one file or folder is supported at the moment.")
		waitForAnyKey()
		return
	}

	// Check installation status
	tools, err := install.GetToolsPath()
	if err != nil {
		pterm.Error.Printfln("Error getting tools path: %s", err)
		waitForAnyKey()
		return
	}

	opts := convert.ConvertOptions{
		Distance: 0.5,
	}

	pterm.DefaultHeader.Println("Settings")
	pterm.Info.Printfln("Exiftool path:   %s", tools.Exiftool)
	pterm.Info.Printfln("cjpegli path:    %s", tools.Cjpegli)
	pterm.Info.Printfln("Jpegli Distance: %.2f (recommended 0.5-3.0)", opts.Distance)

	pterm.DefaultHeader.Println("Converting")
	// Get only JPEG files

	warn := func(msg string) {
		pterm.Warning.Printfln(msg)
	}
	filter := func(path string) bool {
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".jpg", ".jpeg", ".jxl", ".ppm", ".pnm", ".pfm", ".pam", ".pgx", ".png", ".apng", ".gif":
			return true
		default:
			return false
		}
	}

	files, err := filehandling.GetAllFilesInDirectory(filter, os.Args[1], warn)

	if err != nil {
		pterm.Error.Printfln("Error getting files: %s", err)
		waitForAnyKey()
		return
	}
	if len(files) == 0 {
		pterm.Error.Printfln("No compatible image files found in the specified path.")
		pterm.Info.Printfln("Compatible formats: .jpg, .jpeg, .jxl, .ppm, .pnm, .pfm, .pam, .pgx, .png, .apng, .gif")
		waitForAnyKey()
		return
	}

	isDir, err := filehandling.IsPathDir(os.Args[1])
	if err != nil {
		pterm.Error.Printfln("Error checking if path is a directory: %s", err)
		waitForAnyKey()
		return
	}

	states := []convert.ConvertStats{}

	if !isDir {
		stat, err := convert.Convert(tools, opts, files[0], files[0]+".jpegli.jpg")
		if err != nil {
			pterm.Error.Printfln("Error converting file: %s", err)
			waitForAnyKey()
			return
		}
		states = append(states, stat)
		pterm.Info.Printfln("Converted file: %s with ratio %.2f", files[0], stat.FileSizeRatio)
	} else {
		targetFolder := os.Args[1] + "_jpegli-optimized"
		err := os.MkdirAll(targetFolder, os.ModePerm)
		if err != nil {
			pterm.Error.Printfln("Error creating target folder: %s", err)
			waitForAnyKey()
			return
		}
		p, _ := pterm.DefaultProgressbar.WithTotal(len(files)).WithTitle("Converting files").Start()
		for _, file := range files {
			p.UpdateTitle(fmt.Sprintf("Converting %s", file))
			baseName := filepath.Base(file)
			ext := strings.ToLower(filepath.Ext(baseName))
			if ext != ".jpg" && ext != ".jpeg" {
				baseName = strings.TrimSuffix(baseName, ext) + ".jpg"
			}
			targetFilePath := targetFolder + string(os.PathSeparator) + baseName
			stat, err := convert.Convert(tools, opts, file, targetFilePath)
			if err != nil {
				pterm.Error.Printfln("Error converting file: %s", err)
				waitForAnyKey()
				return
			} else {
				states = append(states, stat)
				pterm.Info.Printfln("Converted file: %s with ratio %.2f", file, stat.FileSizeRatio)
			}
			p.Increment()
		}
		p.Stop()
		pterm.Info.Printfln("Converted %d files to %s", len(files), targetFolder)
	}

	pterm.DefaultHeader.Println("Finished")

	// Print the conversion statistics
	var totalSourceSize int64
	var totalTargetSize int64
	for _, stat := range states {
		totalSourceSize += stat.SourceSize
		totalTargetSize += stat.TargetSize
	}

	savedSpace := float64(totalSourceSize-totalTargetSize) / (1024 * 1024) // Convert to MB
	pterm.Success.Printfln("Total space saved: %.2f MB", savedSpace)
	pterm.Info.Printfln("Original size: %.2f MB, New size: %.2f MB",
		float64(totalSourceSize)/(1024*1024),
		float64(totalTargetSize)/(1024*1024))
	pterm.Info.Printfln("Average compression ratio: %.2f%%",
		(1-float64(totalTargetSize)/float64(totalSourceSize))*100)

	waitForAnyKey()
}

func boolToText(b bool) string {
	if b {
		return pterm.Green("Yes")
	}
	return pterm.Red("No")
}
