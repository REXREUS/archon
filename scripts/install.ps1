$Repo = "rexreus/archon" # Replace with the correct GitHub repository
$AppName = "archon"

Write-Host "Installing $AppName from GitHub Release..." -ForegroundColor Cyan

# Detect Architecture
$Arch = if ($Is64Bit) { "amd64" } else { "386" }
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { $Arch = "arm64" }

$OS = "windows"
$BinaryName = "$AppName-$OS-$Arch.exe"

# Fetch latest version from GitHub API
try {
    $ReleaseInfo = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
    $LatestTag = $ReleaseInfo.tag_name
} catch {
    Write-Host "Failed to fetch latest version from GitHub. Ensure REPO '$Repo' is correct and has a release." -ForegroundColor Red
    return
}

$DownloadUrl = "https://github.com/$Repo/releases/download/$LatestTag/$BinaryName"

Write-Host "Downloading $BinaryName version $LatestTag..."
try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile "$BinaryName" -ErrorAction Stop
} catch {
    Write-Host "Download failed. Ensure release URL is correct: $DownloadUrl" -ForegroundColor Red
    return
}

if (Test-Path "$BinaryName") {
    Write-Host "Download successful." -ForegroundColor Green
    
    # Suggestion to move to a folder in PATH
    Write-Host "Suggestion: Move $BinaryName to a folder in your PATH and rename it to $AppName.exe"
    Write-Host "Example: Move-Item $BinaryName C:\Windows\System32\archon.exe"
} else {
    Write-Host "Download failed." -ForegroundColor Red
}
