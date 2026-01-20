# Configuration
$Repo = "rexreus/archon"
$AppName = "archon"
$InstallDir = Join-Path $HOME ".$AppName\bin"

# 1. Security & Environment Setup
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
Write-Host "--- Starting $AppName Installation ---" -ForegroundColor Cyan

# 2. Robust Architecture Detection
$Is64Bit = [Environment]::Is64BitOperatingSystem
$ProcessorArch = $env:PROCESSOR_ARCHITECTURE
$Arch = if ($ProcessorArch -eq "ARM64") { "arm64" } elseif ($Is64Bit) { "amd64" } else { "386" }
$OS = "windows"
$BinaryName = "$AppName-$OS-$Arch.exe"

# 3. Create Installation Directory
if (!(Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    Write-Host "Created installation directory: $InstallDir" -ForegroundColor Gray
}

# 4. Fetch Latest Release Information
try {
    $ApiUrl = "https://api.github.com/repos/$Repo/releases/latest"
    $ReleaseInfo = Invoke-RestMethod -Uri $ApiUrl -Method Get -Headers @{"User-Agent"="PowerShell-Install-Script"}
    $LatestTag = $ReleaseInfo.tag_name
} catch {
    Write-Host "CRITICAL: Could not reach GitHub API. Check your internet connection." -ForegroundColor Red
    return
}

$DownloadUrl = "https://github.com/$Repo/releases/download/$LatestTag/$BinaryName"
$FinalExecutable = Join-Path $InstallDir "$AppName.exe"

# 5. Download and Deploy
$ProgressPreference = 'SilentlyContinue'
Write-Host "Downloading $AppName $LatestTag for $Arch..." -ForegroundColor Yellow
try {
    # Download to temporary location first to ensure integrity
    $TempFile = Join-Path $env:TEMP "$BinaryName"
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempFile -ErrorAction Stop
    
    # Move and Rename to final location
    Move-Item -Path $TempFile -Destination $FinalExecutable -Force
} catch {
    Write-Host "ERROR: Download failed. The release asset might not exist: $DownloadUrl" -ForegroundColor Red
    return
} finally {
    $ProgressPreference = 'Continue'
}

# 6. Persistent PATH Integration
Write-Host "Updating Environment PATH..." -ForegroundColor Cyan
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    $NewPath = "$UserPath;$InstallDir"
    [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
    # Update current session path immediately
    $env:PATH += ";$InstallDir"
    Write-Host "Successfully added $InstallDir to User PATH." -ForegroundColor Green
} else {
    Write-Host "PATH already configured." -ForegroundColor Gray
}

# 7. Verification
if (Test-Path $FinalExecutable) {
    Write-Host "`nSUCCESS: $AppName has been installed to $InstallDir" -ForegroundColor Green
    Write-Host "You can now run '$AppName' in any NEW terminal window." -ForegroundColor White
}
