# Remove context menu for all files
Remove-Item -Path "HKCU:\SOFTWARE\Classes\*\shell\JPEGLIOptimizer" -Recurse -ErrorAction SilentlyContinue

# Remove context menu for directories
Remove-Item -Path "HKCU:\SOFTWARE\Classes\Directory\shell\JPEGLIOptimizer" -Recurse -ErrorAction SilentlyContinue

# Remove application data folder
$cacheDir = [System.Environment]::GetFolderPath('LocalApplicationData')
$appFolder = Join-Path $cacheDir 'jpegli-windows-explorer-extension'
if (Test-Path $appFolder) {
    Remove-Item -Path $appFolder -Recurse -Force -ErrorAction SilentlyContinue
    Write-Host "Deleted application data folder: $appFolder"
} else {
    Write-Host "Application data folder not found: $appFolder"
}

Write-Host "JPEGLI Optimizer context menu entries and app data removed."