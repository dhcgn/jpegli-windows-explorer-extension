<#
.SYNOPSIS
    Converts JPEG images using the JPEGLI encoder and preserves metadata.

.DESCRIPTION
    This script processes JPEG images from a source directory, converts them using the JPEGLI encoder
    with specified quality settings, and copies the original metadata to the converted images.
    It also provides size comparison between original and converted files.

.PARAMETER source
    The source directory containing original JPEG images.

.PARAMETER target
    The target directory where converted images will be saved.

.PARAMETER skipIfExists
    If true, skips conversion if the target file already exists.

.NOTES
    Requires:
    - ExifTool (https://exiftool.org/)
    - JPEGLI (https://github.com/google/jpegli)

.EXAMPLE
    .\convert_with_jpegli.ps1
#>

# Check for required executables

# Download from https://exiftool.org/ (need to rename to exiftool.exe, because this changes bahavior)
$exifToolPath = "C:\tools\exiftool-13.26_64\exiftool.exe"
if (-not (Test-Path $exifToolPath)) {
    Write-Error "ExifTool not found at $exifToolPath"
    exit 1
}

# Download from this Githunb Action https://github.com/google/jpegli/actions/workflows/release.yaml to main branch
$cjpegliPath = "C:\tools\jpegli\bin\cjpegli.exe"
if (-not (Test-Path $cjpegliPath)) {
    Write-Error "cjpegli.exe not found at $cjpegliPath"
    exit 1
}

# Define source and target directories
$source = "H:\media\development\2025-04-19\Export"
$target = "H:\media\development\2025-04-19\Export_jpegli"

# skip if target file already exists
$skipIfExists = $true

# Create target folder if it doesn't exist
if (-not (Test-Path $target)) {
    New-Item -ItemType Directory -Path $target | Out-Null
}

# Process the first file from the source folder
# Get-ChildItem $source | Select -First 1 | ForEach-Object {
Get-ChildItem $source -Filter *.jpg | ForEach-Object {
    $sourceFile = $_.FullName
    $targetFile = Join-Path $target $_.Name

    if ($skipIfExists -and (Test-Path $targetFile)) {
        Write-Host "Skipping $targetFile as it already exists."
        return
    }

    # Convert image with cjpegli.exe using -d 1.0 option (1.0 = visually lossless (default)
    & $cjpegliPath $sourceFile $targetFile -d 0.5

    # Copy metadata from source to target using ExifTool
    & $exifToolPath -overwrite_original -TagsFromFile $sourceFile $targetFile

    # Retrieve file sizes
    $s_size = (Get-Item $sourceFile).Length
    $t_size = (Get-Item $targetFile).Length

    # Calculate and display the size ratio
    $ratio = $t_size / $s_size
    Write-Host "Source: $sourceFile ($s_size bytes)"
    Write-Host "Target: $targetFile ($t_size bytes)"
    Write-Host "Size ratio: $([math]::Round($ratio,2))"
}

# Compare the file sizes of the source and target files
$sourceFiles = @(); $sourceFiles += Get-ChildItem $source -Filter *.jpg
$targetFiles = @(); $targetFiles += Get-ChildItem $target -Filter *.jpg
$allFilesGroup = $sourceFiles + $targetFiles | Group-Object -Property Name 
$allFilesGroup = $allFilesGroup | ?{ $_.Count -eq 2 }

# Number of files on both sides
Write-Host "Number of files on both sides: $($allFilesGroup.Count / 2)"

$sum_source = 0
$sum_target = 0

$allFilesGroup | ForEach-Object {
    $sourceFile = $_.Group[0]
    $targetFile = $_.Group[1]

    $s_size = $sourceFile.Length
    $t_size = $targetFile.Length

    $sum_source += $s_size
    $sum_target += $t_size

    # Calculate and display the size ratio
    $ratio = $t_size / $s_size
    Write-Host ("{0} from {1:N2} MB to {2:N2} MB, ratio: {3}" -f ($sourceFile.Name), ($s_size/1mb), ($t_size/1mb), ([math]::Round($ratio,2)))
    # Write-Host "File: $sourceFile ($s_size bytes)"
    # Write-Host "Target: $targetFile ($t_size bytes)"
    # Write-Host "Size ratio: $([math]::Round($ratio,2))"
}

Write-Host ("Total size: {0:N2} MB to {1:N2} MB, ratio: {2}" -f ($sum_source/1mb), ($sum_target/1mb), ([math]::Round($sum_target/$sum_source,2)))

# Format-Table -InputObject $allFiles -Property FullName, Length -AutoSize


