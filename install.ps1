$ErrorActionPreference = 'Stop'

$ExecutableName = "rocketblend"

$APIUrl = "https://api.github.com/repos/rocketblend/rocketblend/releases/latest"

# Fetch the latest release info
$LatestRelease = Invoke-RestMethod -Uri $APIUrl

# Determine the OS and Architecture
$OS = if ([Environment]::Is64BitOperatingSystem) { "Windows_x86_64" } else { "Windows_i386" }

# Find the desired asset based on the OS and Architecture
$DownloadUrl = ($LatestRelease.assets | Where-Object { $_.browser_download_url -like "*$ExecutableName_$OS.zip" }).browser_download_url

if (!$DownloadUrl) {
    throw "Failed to find the download URL for the $OS version of $ExecutableName"
}

# Download and extract the ZIP file
$OutputPath = "$Env:LOCALAPPDATA/Programs/RocketBlend/"
$ZipFile = "$Env:TEMP/$ExecutableName.zip"

Invoke-WebRequest -Uri $DownloadUrl -OutFile $ZipFile
Expand-Archive -Path $ZipFile -DestinationPath $Env:TEMP -Force

# Move the content of the extracted folder to the OutputPath
$ExtractedFolderPath = "$Env:TEMP/${ExecutableName}_${OS}"

# Ensure the OutputPath directory exists
if (-not (Test-Path -Path $OutputPath -PathType Container)) {
    New-Item -Path $OutputPath -ItemType Directory | Out-Null
}

Move-Item -Path "$ExtractedFolderPath/*" -Destination $OutputPath -Force

# Cleanup
Remove-Item $ZipFile
Remove-Item $ExtractedFolderPath -Recurse

# Add the installation path to the PATH environment variable
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
$PathArray = $UserPath -split ";"
if (-not ($PathArray -contains $OutputPath)) {
    $NewUserPath = $UserPath + ";" + $OutputPath
    [Environment]::SetEnvironmentVariable("Path", $NewUserPath, "User")
}

Write-Host "$ExecutableName has been installed successfully!"
Write-Host "Please restart your PowerShell session for the changes to take effect."