//go:build windows
// +build windows

package install

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

// SetExecutableAsWindowsExplorerContextMenu sets the executable as a Windows Explorer context menu item
// for Files and Folders.
func SetExecutableAsWindowsExplorerContextMenu(execPath string) {
	// Keep the path as is (don't convert to slashes) and properly escape it for the registry
	execCommand := "\"" + execPath + "\" \"%1\""
	// Registry keys to modify
	registryKeys := []struct {
		parent string
		path   string
		name   string
		value  string
	}{
		// For all files
		{`SOFTWARE\Classes\*\shell`, "JPEGLIOptimizer", "", "Optimize with JPEGLI"},
		{`SOFTWARE\Classes\*\shell\JPEGLIOptimizer`, "Icon", "", execPath},
		{`SOFTWARE\Classes\*\shell\JPEGLIOptimizer\command`, "", "", execCommand},

		// For folders
		{`SOFTWARE\Classes\Directory\shell`, "JPEGLIOptimizer", "", "Optimize JPEGs with JPEGLI"},
		{`SOFTWARE\Classes\Directory\shell\JPEGLIOptimizer`, "Icon", "", execPath},
		{`SOFTWARE\Classes\Directory\shell\JPEGLIOptimizer\command`, "", "", execCommand},
	}

	// Create or update registry entries
	for _, key := range registryKeys {
		// Try to open existing key first
		fullPath := key.parent + "\\" + key.path
		k, exists, err := registry.CreateKey(registry.CURRENT_USER, fullPath, registry.ALL_ACCESS)
		if err != nil {
			fmt.Printf("Error accessing registry key %s: %v\n", fullPath, err)
			continue
		}

		// Set or update the registry value
		if key.name == "" {
			// Default value
			err = k.SetStringValue("", key.value)
			if err != nil {
				fmt.Printf("Error setting default value for key %s: %v\n", fullPath, err)
			} else if exists {
				fmt.Printf("Updated existing registry key: %s\n", fullPath)
			} else {
				fmt.Printf("Created new registry key: %s\n", fullPath)
			}
		} else {
			// Named value
			err = k.SetStringValue(key.name, key.value)
			if err != nil {
				fmt.Printf("Error setting value '%s' for key %s: %v\n", key.name, fullPath, err)
			} else if exists {
				fmt.Printf("Updated existing registry value: %s\\%s\n", fullPath, key.name)
			} else {
				fmt.Printf("Created new registry value: %s\\%s\n", fullPath, key.name)
			}
		}

		k.Close()
	}

	fmt.Println("Successfully updated JPEGLI Optimizer in Windows Explorer context menu")
}
