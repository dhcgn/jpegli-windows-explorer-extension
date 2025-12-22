package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	update "github.com/dhcgn/gh-update"
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

const (
	ExitCodeSuccess         = 0
	ExitCodeSettingsError   = 1
	ExitCodeToolsMissing    = 2
	ExitCodePathError       = 3
	ExitCodeNoFiles         = 4
	ExitCodeConversionError = 5
)

var (
	Version = "UNSET"
	Build   = "UNSET"
	Commit  = "UNSET"
)

type App struct {
	NoUserInteraction bool
}

// pauseInConsole simulates a pause in console by printing dots
func pauseInConsole() {
	for i := 0; i < 3; i++ {
		fmt.Print(".")
		time.Sleep(1 * time.Second)
	}
	fmt.Println()
}

func (a *App) WaitForAnyKey() {
	if a.NoUserInteraction {
		pauseInConsole()
		return
	}

	fmt.Println("Press any key to continue...")
	var input [1]byte
	os.Stdin.Read(input[:])
}

func main() {
	os.Exit(Run(os.Args, nil))
}

func Run(args []string, opts *settings.Seetings) int {
	app := &App{}

	fmt.Println("jpegli-windows-explorer-extension")
	fmt.Printf("Version: %s, Build: %s, Commit: %s\n", Version, Build, Commit)

	// if one of the args is --help or -h, show help and exit
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			pterm.Println("Usage: jpegli-windows-explorer-extension [file1 file2 ... | directory]")
			pterm.Println("Documentation: https://github.com/dhcgn/jpegli-windows-explorer-extension/blob/main/README.md")
			app.WaitForAnyKey()
			return ExitCodeSuccess
		}
	}

	var cfgPath string
	if opts == nil {
		if !settings.CheckForConfigFile() {
			pterm.Warning.Println("No configuration file found, creating default configuration at " + settings.GetConfigFilePath())
		}

		loadedOpts, path, err := settings.LoadOrDefault()
		if err != nil {
			pterm.Warning.Printfln("Error loading settings, using defaults: %s", err)
			app.WaitForAnyKey()
			return ExitCodeSettingsError
		}
		opts = &loadedOpts
		cfgPath = path
	} else {
		cfgPath = "provided"
	}

	app.NoUserInteraction = opts.NoUserInteraction

	if opts.SkipUpdateCheck {
		pterm.Info.Println("Skipping update check as per configuration.")
	} else {
		pterm.Print("Checking for updates ... ")
		lr, err := update.GetLatestVersion("dhcgn/jpegli-windows-explorer-extension", Version, "^jpegli-windows-explorer-extension.exe$")
		if err == update.ErrorNoNewVersionFound {
			pterm.Info.Println("You are running the latest version.")
		} else if err != nil {
			pterm.Warning.Println("Failed to check for updates.")
			pterm.Warning.Printfln("Error: %s", err)
			pauseInConsole()
		} else {
			pterm.Printf("New Version: '%s' is available! You have '%s'\n", lr.Version, Version)
			pauseInConsole()
		}
	}

	// If no arguments provided, show install prompt
	if len(args) == 1 {
		handleInstallPrompt()
		app.WaitForAnyKey()
		return ExitCodeSuccess
	}

	tools := getToolsOrExit()
	if tools == nil {
		app.WaitForAnyKey()
		return ExitCodeToolsMissing
	}

	showSettings(tools, *opts, cfgPath)

	if opts.NoUserInteraction {
		pterm.Info.Println("No user interaction mode enabled for processing files.")
	}

	printArgs(args)
	filesOrDirs := args[1:]

	isDir, err := checkIsDirOrExit(filesOrDirs)
	if err != nil {
		pterm.Error.Printfln("Error checking if path is a directory: %s", err)
		app.WaitForAnyKey()
		return ExitCodePathError
	}
	if isDir == nil {
		app.WaitForAnyKey()
		return ExitCodePathError
	}

	files := getFilesOrExit(filesOrDirs)
	if files == nil {
		app.WaitForAnyKey()
		return ExitCodeNoFiles
	}

	var targetDirBase string
	if len(filesOrDirs) > 0 {
		targetDirBase = filesOrDirs[0]
	}

	states := convertFilesOrExit(files, *isDir, tools, *opts, targetDirBase)
	if states == nil {
		app.WaitForAnyKey()
		return ExitCodeConversionError
	}

	printStats(states)
	app.WaitForAnyKey()
	return ExitCodeSuccess
}

func handleInstallPrompt() {
	pterm.Println("No arguments provided. Want to install and set context menu? --help for more info.")
	result, _ := pterm.DefaultInteractiveConfirm.Show()
	pterm.Println()
	pterm.Info.Printfln("You answered: %s", boolToText(result))
	if result {
		install.Do()
		pterm.Println("Installation completed.")
	} else {
		pterm.Println("You chose not to install.")
	}
}

func printArgs(args []string) {
	for i := 1; i < len(args); i++ {
		fmt.Printf("file/folder: %d: %s\n", i, args[i])
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

func showSettings(tools *types.ExecutablePaths, opts settings.Seetings, cfgPath string) {
	pterm.DefaultHeader.Println("Settings")
	pterm.Info.Printfln("Config file:     %s", cfgPath)
	pterm.Info.Printfln("Exiftool path:   %s", tools.Exiftool)
	pterm.Info.Printfln("cjpegli path:    %s", tools.Cjpegli)
	pterm.Info.Printfln("Jpegli Distance: %.2f (recommended 0.5-3.0, 1.0 = visually lossless, lower better)", opts.Distance)
	pterm.Info.Printfln("Override Original: %v", opts.OverrideOriginalFile)
	pterm.DefaultHeader.Println("Converting")
}

func getFilesOrExit(filesOrDirs []string) []string {
	warn := func(msg string) { pterm.Warning.Printfln("%s", msg) }
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

func convertFilesOrExit(files []string, isDir bool, tools *types.ExecutablePaths, opts settings.Seetings, targetDirBase string) []convert.ConvertStats {
	states := []convert.ConvertStats{}
	if !isDir {
		for _, file := range files {
			var targetPath string
			if opts.OverrideOriginalFile {
				// When overriding, use the source file as the target
				targetPath = file
			} else {
				// When not overriding, create a new file with .jpegli.jpg suffix
				baseName := filepath.Base(file)
				ext := filepath.Ext(baseName)
				targetName := strings.TrimSuffix(baseName, ext) + ".jpegli.jpg"
				targetPath = filepath.Join(filepath.Dir(file), targetName)
			}
			stat, err := convert.Convert(*tools, opts.Distance, opts.OverrideOriginalFile, file, targetPath)
			if err != nil {
				pterm.Error.Printfln("Error converting file: %s", err)
				return nil
			}
			states = append(states, stat)
			pterm.Info.Printfln("Converted file: %s with ratio %.2f", file, stat.FileSizeRatio)
		}
	} else {
		targetFolder := targetDirBase + "_jpegli-optimized"
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
			stat, err := convert.Convert(*tools, opts.Distance, opts.OverrideOriginalFile, file, targetFilePath)
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
