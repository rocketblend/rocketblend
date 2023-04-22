$Repo = "https://github.com/rocketblend/rocketblend"
$AppName = "rktb"
$ReleasesApi = "https://api.github.com/repos/$Repo/releases/latest"

# Determine the platform and architecture
$OS = "windows"
$Arch = "amd64" # Currently assuming amd64, modify if needed for other architectures

# Download the appropriate binary for the platform and architecture
$DownloadUrl = (Invoke-WebRequest -Uri $ReleasesApi | ConvertFrom-Json | `
    Where-Object { $_.browser_download_url -match "${OS}-${Arch}" }).browser_download_url
$Destination = Join-Path -Path $HOME -ChildPath ".local\bin\$AppName.exe"

if (-Not (Test-Path -Path (Split-Path -Path $Destination -Parent))) {
    New-Item -ItemType Directory -Path (Split-Path -Path $Destination -Parent) | Out-Null
}

Write-Host "Downloading $AppName for $OS-$Arch..."
Invoke-WebRequest -Uri $DownloadUrl -OutFile $Destination

# Add the destination directory to the user's PATH, if not already present
$PathKey = [Environment]::GetEnvironmentVariable("Path", "User")
$DestinationDir = Split-Path -Path $Destination -Parent

if (-Not ($PathKey -split ";" -contains $DestinationDir)) {
    Write-Host "Adding $AppName to the PATH..."
    [Environment]::SetEnvironmentVariable("Path", $PathKey + ";" + $DestinationDir, "User")
}

Write-Host "Installation complete. $AppName is now available in $Destination"