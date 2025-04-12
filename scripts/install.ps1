# Default version
$DefaultVersion = "v0.0.1"

# Get the version from the command line argument or use default
$Version = $args[0]
if (-not $Version) {
    $Version = $DefaultVersion
}

# Define the base URL for the release artifacts
$BaseUrl = "https://github.com/struckchure/gv/releases/download/$Version"

# Define the file names (adjust these as needed)
$Windows_AMD64 = "gv_Windows_x86_64.zip"
$Windows_ARM64 = "gv_Windows_arm64.zip"

# Determine the architecture
$Arch = [System.Environment]::GetEnvironmentVariable("PROCESSOR_ARCHITECTURE")

# Set the file to download based on architecture
switch ($Arch) {
    "AMD64" {
        $File = $Windows_AMD64
    }
    "ARM64" {
        $File = $Windows_ARM64
    }
    default {
        Write-Host "Unsupported architecture: $Arch"
        exit 1
    }
}

# Define the destination directory
$DestDir = "$env:USERPROFILE\.gv\bin"

# Create the destination directory if it does not exist
if (-not (Test-Path $DestDir)) {
    New-Item -Path $DestDir -ItemType Directory | Out-Null
}

# Download the file
Write-Host "Downloading $File..."
Invoke-WebRequest -Uri "$BaseUrl/$File" -OutFile "$File"

# Extract the downloaded file to the .gv\bin directory
Write-Host "Extracting $File to $DestDir..."
Expand-Archive -Path $File -DestinationPath $DestDir -Force

# Remove the downloaded file
Remove-Item -Path $File -Force

# Add .gv\bin to the PATH if it's not already there
$PathEntry = "$DestDir"
$CurrentPath = [System.Environment]::GetEnvironmentVariable("PATH", [System.EnvironmentVariableTarget]::User)

if ($CurrentPath -notlike "*$PathEntry*") {
    [System.Environment]::SetEnvironmentVariable("PATH", "$CurrentPath;$PathEntry", [System.EnvironmentVariableTarget]::User)
    Write-Host "Updated PATH to include $DestDir"
} else {
    Write-Host "$DestDir is already in PATH"
}
