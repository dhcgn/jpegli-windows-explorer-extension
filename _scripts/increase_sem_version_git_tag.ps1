# increase the sem version with git tag

# Get the current version
$version = git describe --tags --abbrev=0 2>$null

# If no tag exists, start with v0.0.1
if (-not $version) {
    $version = "v0.0.1"
} else {
    # Extract version numbers
    if ($version -match 'v?(\d+)\.(\d+)\.(\d+)') {
        $major = [int]$matches[1]
        $minor = [int]$matches[2]
        $patch = [int]$matches[3]
        
        # Increment patch version
        $patch++
        
        # Create new version string
        $version = "v$major.$minor.$patch"
    } else {
        Write-Error "Current tag is not in semantic version format"
        exit 1
    }
}

# Create the new tag
git tag $version

# Push the new tag
Write-Host "New version tag: $version"
Write-Host "Execute: " -NoNewline
Write-Host "git push origin $version" -ForegroundColor Magenta