# Configuration
$Repo = "rexreus/archon"
$AppName = "archon"

# 1. Security: Ensure TLS 1.2 is used for GitHub API/Download (Crucial for PowerShell 5.1 compatibility)
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

Write-Host "Installing $AppName from GitHub Release..." -ForegroundColor Cyan

# 2. Robust Architecture Detection
# We use .NET Framework methods to ensure it works on both Windows PowerShell 5.1 and PowerShell 7+
$Is64Bit = [Environment]::Is64BitOperatingSystem
$ProcessorArch = $env:PROCESSOR_ARCHITECTURE

$Arch = if ($ProcessorArch -eq "ARM64") { 
    "arm64" 
} elseif ($Is64Bit) { 
    "amd64" 
} else { 
    "386" 
}

$OS = "windows"
$BinaryName = "$AppName-$OS-$Arch.exe"

# 3. Fetch Latest Release Information
try {
    $ApiUrl = "https://api.github.com/repos/$Repo/releases/latest"
    # Using -UseBasicParsing for speed and compatibility
    $ReleaseInfo = Invoke-RestMethod -Uri $ApiUrl -Method Get
    $LatestTag = $ReleaseInfo.tag_name
} catch {
    Write-Host "ERROR: Failed to fetch latest version from GitHub API." -ForegroundColor Red
    Write-Host "Please verify if the repository '$Repo' exists and has a public release." -ForegroundColor Gray
    return
}

$DownloadUrl = "https://github.com/$Repo/releases/download/$LatestTag/$BinaryName"

# 4. Optimized Download Process
# Temporarily disabling progress bar significantly increases download speed in CLI environments
$OldProgressPreference = $ProgressPreference
$ProgressPreference = 'SilentlyContinue'

Write-Host "Downloading version $LatestTag ($Arch)..." -ForegroundColor Yellow
try {
    $OutputPath = Join-Path $pwd $BinaryName
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $OutputPath -ErrorAction Stop
} catch {
    Write-Host "ERROR: Download failed. File not found at: $DownloadUrl" -ForegroundColor Red
    Write-Host "Note: Ensure the asset naming convention matches '$BinaryName' in the GitHub Release." -ForegroundColor Gray
    return
} finally {
    # Restore Progress Bar setting
    $ProgressPreference = $OldProgressPreference
}

# 5. Verification and Post-Install Suggestion
if (Test-Path $OutputPath) {
    $FileSize = (Get-Item $OutputPath).Length / 1MB
    Write-Host "Download successful! (Size: $('{0:N2}' -f $FileSize) MB)" -ForegroundColor Green
    
    Write-Host "`nInstallation Complete." -ForegroundColor White
    Write-Host "Next Step: Move the binary to a folder in your System PATH to use it globally." -ForegroundColor Gray
    Write-Host "Example: Move-Item `"$BinaryName`" `"$env:ProgramFiles\$AppName.exe`" -Force" -ForegroundColor DarkGray
} else {
    Write-Host "ERROR: File verification failed. The binary was not saved correctly." -ForegroundColor Red
}
