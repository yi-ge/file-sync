Write-Host "Detecting operating system and architecture..."

$OS = (Get-WmiObject -Class Win32_OperatingSystem).Caption.ToLower()
$ARCH = (Get-WmiObject -Class Win32_Processor).AddressWidth

switch ($ARCH) {
    64 { $ARCH = "amd64" }
    32 { $ARCH = "386" }
    default {
        Write-Host "Unsupported architecture: $ARCH"
        exit 1
    }
}

if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Host "This script must be run with administrator privileges on Windows."
    Write-Host "Right-click on the terminal and select 'Run as administrator', then try again."
    exit 1
}

$TARGET_PATH = "$env:SYSTEMROOT\System32"

Write-Host "Fetching the latest version number..."
$VERSION = (Invoke-RestMethod -Uri "https://api.github.com/repos/yi-ge/file-sync/releases/latest").tag_name
if ([string]::IsNullOrEmpty($VERSION)) {
    Write-Host "Failed to fetch the latest version number"
    exit 1
}

$FILENAME = "file-sync-${OS}-${ARCH}.exe"

Write-Host "Downloading $FILENAME version $VERSION..."

$DOWNLOAD_URL = "https://github.com/yi-ge/file-sync/releases/download/${VERSION}/${FILENAME}"
Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $FILENAME -ErrorAction Stop -ErrorVariable DOWNLOAD_ERROR

if ($DOWNLOAD_ERROR) {
    Write-Host "Download failed"
    exit 1
}

Write-Host "Moving ${FILENAME} to ${TARGET_PATH}..."

Move-Item -Path $FILENAME -Destination "${TARGET_PATH}\file-sync" -Force

Write-Host "Checking if the file-sync is working properly..."
try {
    & "${TARGET_PATH}\file-sync" -v
} catch {
    Write-Host "The file-sync command does not work as expected"
    exit 1
}

Write-Host "Registering file-sync as a service..."
& "${TARGET_PATH}\file-sync" service enable

Write-Host "Starting file-sync service..."
& "${TARGET_PATH}\file-sync" service start

Write-Host "Done! File-sync has been successfully set up, registered, and started as a service."
