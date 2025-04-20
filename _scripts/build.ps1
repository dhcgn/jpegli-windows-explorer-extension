# Set build variables
$VERSION = git describe --tags --always --dirty 2>$null
if (-not $VERSION) { $VERSION = 'dev' }
$BUILD = Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ'
$COMMIT = git rev-parse --short HEAD

# Run tests
Write-Host 'Running tests...'
go test -v ./...

# Build with ldflags
Write-Host 'Building executable...'
go build -o jpegli-windows-explorer-extension.exe -ldflags "-X 'main.Version=$VERSION' -X 'main.Build=$BUILD' -X 'main.Commit=$COMMIT'" .\main.go

Write-Host 'Build complete. Output: jpegli-windows-explorer-extension.exe'