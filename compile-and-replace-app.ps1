go build -o app.exe .\main.go 
$hash = Get-FileHash .\app.exe -Algorithm SHA256
$installedExePath = Join-Path $env:LOCALAPPDATA "jpegli-windows-explorer-extension\jpegli-windows-explorer-extension.exe"

if (Test-Path $installedExePath) {
    $installedHash = Get-FileHash $installedExePath -Algorithm SHA256
    if ($hash.Hash -eq $installedHash.Hash) {
        Write-Host "The built app.exe is identical to the installed version. No replacement needed."
    } else {
        Write-Host "The built app.exe differs from the installed version. Replacing..."
    }
} else {
    Write-Host "No installed version found. Installing the built app.exe..."
}

Copy-Item -Path .\app.exe -Destination $installedExePath -Force

Write-Host "Replacement complete."

$configPath = Join-Path $env:LOCALAPPDATA "jpegli-windows-explorer-extension\config.yaml"
if (Test-Path $configPath) {
    Write-Host "Preserving existing configuration file at $configPath"
    Write-Host "-----------"
    Get-Content $configPath | ForEach-Object { Write-Host $_ }
    Write-Host "-----------"
} else {
    Write-Host "No existing configuration file found. "
}
