package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhcgn/jpegli-windows-explorer-extension/convert"
	"github.com/dhcgn/jpegli-windows-explorer-extension/filehandling"
	"github.com/dhcgn/jpegli-windows-explorer-extension/install"
	"github.com/dhcgn/jpegli-windows-explorer-extension/settings"
	"github.com/dhcgn/jpegli-windows-explorer-extension/types"
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

	if handleInstallPrompt() {
		waitForAnyKey()
		return
	}

	tools := getToolsOrExit()
	if tools == nil {
		waitForAnyKey()
		return
	}

	opts, err := settings.LoadOrDefault()
	if err != nil {
		pterm.Warning.Printfln("Error loading settings, using defaults: %s", err)
		waitForAnyKey()
		return
	}
	showSettings(tools, opts)

	printArgs()
	filesOrDirs := os.Args[1:]

	isDir, err := checkIsDirOrExit(filesOrDirs)
	if err != nil {
		pterm.Error.Printfln("Error checking if path is a directory: %s", err)
		waitForAnyKey()
		return
	}
	if isDir == nil {
		waitForAnyKey()
		return
	}

	files := getFilesOrExit(filesOrDirs)
	if files == nil {
		waitForAnyKey()
		return
	}

	states := convertFilesOrExit(files, *isDir, tools, opts)
	if states == nil {
		waitForAnyKey()
		return
	}

	printStats(states)
	waitForAnyKey()
}

func handleInstallPrompt() bool {
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
		return true
	}
	return false
}

func printArgs() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("file/folder: %d: %s\n", i, os.Args[i])
	}
}

func getToolsOrExit() *types.ExecutablePaths {
	tools, err := install.GetToolsPath()
	if err != nil {
		pterm.Error.Printfln("Error getting tools path: %s", err)
		return nil
	}
	return &tools
}

func showSettings(tools *types.ExecutablePaths, opts convert.ConvertOptions) {
	pterm.DefaultHeader.Println("Settings")
	pterm.Info.Printfln("Exiftool path:   %s", tools.Exiftool)
	pterm.Info.Printfln("cjpegli path:    %s", tools.Cjpegli)
	pterm.Info.Printfln("Jpegli Distance: %.2f (recommended 0.5-3.0, 1.0 = visually lossless, lower better)", opts.Distance)
	pterm.DefaultHeader.Println("Converting")
}

func getFilesOrExit(filesOrDirs []string) []string {
	warn := func(msg string) { pterm.Warning.Printfln(msg) }
	filter := func(path string) bool {
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".jpg", ".jpeg", ".jxl", ".ppm", ".pnm", ".pfm", ".pam", ".pgx", ".png", ".apng", ".gif":
			return true
		default:
			return false
		}
	}

	var files []string
	for _, path := range filesOrDirs {
		moreFiles, err := filehandling.GetAllFilesInDirectory(filter, path, warn)
		if err != nil {
			pterm.Error.Printfln("Error getting files: %s", err)
			return nil
		}
		files = append(files, moreFiles...)
	}

	if len(files) == 0 {
		pterm.Error.Printfln("No compatible image files found in the specified path.")
		pterm.Info.Printfln("Compatible formats: .jpg, .jpeg, .jxl, .ppm, .pnm, .pfm, .pam, .pgx, .png, .apng, .gif")
		return nil
	}
	return files
}

func checkIsDirOrExit(filesOrDirs []string) (*bool, error) {
	if len(filesOrDirs) == 0 {
		return nil, fmt.Errorf("no files or directories provided")
	}

	// If one element and it's a directory, return true
	if len(filesOrDirs) == 1 {
		isDir, err := filehandling.IsPathDir(filesOrDirs[0])
		if err != nil {
			return nil, err
		}
		if isDir {
			result := true
			return &result, nil
		}
	}

	// Check if all elements are files
	allFiles := true
	for _, path := range filesOrDirs {
		isDir, err := filehandling.IsPathDir(path)
		if err != nil {
			return nil, err
		}
		if isDir {
			allFiles = false
			break
		}
	}

	if allFiles {
		result := false
		return &result, nil
	}

	// All other cases return error
	return nil, fmt.Errorf("invalid combination: must be either a single directory or multiple files only")
}

func convertFilesOrExit(files []string, isDir bool, tools *types.ExecutablePaths, opts convert.ConvertOptions) []convert.ConvertStats {
	states := []convert.ConvertStats{}
	if !isDir {
		for _, file := range files {
			baseName := filepath.Base(file)
			ext := filepath.Ext(baseName)
			targetName := strings.TrimSuffix(baseName, ext) + ".jpegli.jpg"
			targetPath := filepath.Join(filepath.Dir(file), targetName)
			stat, err := convert.Convert(*tools, opts, file, targetPath)
			if err != nil {
				pterm.Error.Printfln("Error converting file: %s", err)
				return nil
			}
			states = append(states, stat)
			pterm.Info.Printfln("Converted file: %s with ratio %.2f", file, stat.FileSizeRatio)
		}
	} else {
		targetFolder := os.Args[1] + "_jpegli-optimized"
		err := os.MkdirAll(targetFolder, os.ModePerm)
		if err != nil {
			pterm.Error.Printfln("Error creating target folder: %s", err)
			return nil
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
			stat, err := convert.Convert(*tools, opts, file, targetFilePath)
			if err != nil {
				pterm.Error.Printfln("Error converting file: %s", err)
				return nil
			} else {
				states = append(states, stat)
				pterm.Info.Printfln("Converted file: %s with ratio %.2f", file, stat.FileSizeRatio)
			}
			p.Increment()
		}
		p.Stop()
		pterm.Info.Printfln("Converted %d files to %s", len(files), targetFolder)
	}
	return states
}

func printStats(states []convert.ConvertStats) {
	pterm.DefaultHeader.Println("Finished")
	var totalSourceSize int64
	var totalTargetSize int64
	for _, stat := range states {
		totalSourceSize += stat.SourceSize
		totalTargetSize += stat.TargetSize
	}
	savedSpace := float64(totalSourceSize-totalTargetSize) / (1024 * 1024)
	pterm.Success.Printfln("Total space saved: %.2f MB", savedSpace)
	pterm.Info.Printfln("Original size: %.2f MB, New size: %.2f MB",
		float64(totalSourceSize)/(1024*1024),
		float64(totalTargetSize)/(1024*1024))
	pterm.Info.Printfln("Average compression ratio: %.2f%%",
		(1-float64(totalTargetSize)/float64(totalSourceSize))*100)
}

func boolToText(b bool) string {
	if b {
		return pterm.Green("Yes")
	}
	return pterm.Red("No")
}
