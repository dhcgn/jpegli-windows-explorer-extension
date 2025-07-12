$inputFile = gci "test-files\DSC_4045-NEF_DxO_DeepPRIME.jpg"
go run .\main.go $inputFile.FullName

<#

Expected Output:

jpegli-windows-explorer-extension
Version: UNSET, Build: UNSET, Commit: UNSET
file/folder: 1: test-files\DSC_4045-NEF_DxO_DeepPRIME.jpg
                  
     Settings     
                  
 INFO  Exiftool path:   C:\Users\USER\AppData\Local\jpegli-windows-explorer-extension\exiftool-13.27_64\exiftool-13.27_64\exiftool.exe
 INFO  cjpegli path:    C:\Users\USER\AppData\Local\jpegli-windows-explorer-extension\jpegli-x64-windows-static\bin\cjpegli.exe
 INFO  Jpegli Distance: 0.50 (recommended 0.5-3.0, 1.0 = visually lossless, lower better)
                    
     Converting     
                    
 INFO  Converted file: test-files\DSC_4045-NEF_DxO_DeepPRIME.jpg with ratio 0.60

     Finished     
                  
 SUCCESS  Total space saved: 9.10 MB
 INFO  Original size: 22.51 MB, New size: 13.41 MB
 INFO  Average compression ratio: 40.42%
#>

Write-Host "Start Comparison"

$exe = Join-Path $env:LOCALAPPDATA "Programs\Beyond Compare 5\BCompare.exe"
$args = @($inputFile.FullName, (Join-Path $inputFile.Directory  ($inputFile.BaseName + ".jpegli.jpg")))

Write-Host "Running:   $exe"
Write-Host "With args: $args"

if (-not (Test-Path $exe)) {
    Write-Host "Beyond Compare executable not found at $exe"
    exit 1
}

& $exe $args