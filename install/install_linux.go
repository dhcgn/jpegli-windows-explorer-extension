//go:build linux
// +build linux

package install

// SetExecutableAsWindowsExplorerContextMenu sets the executable as a Windows Explorer context menu item
// for Files and Folders.
func SetExecutableAsWindowsExplorerContextMenu(execPath string) {
	// This function is a no-op on Linux, as context menu integration is handled differently.
	// You can implement this function if you want to add context menu integration for Linux.
}
