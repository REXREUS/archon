# Configuration
$AppName = "archon"
$BinaryName = "$AppName.exe"
$InstallDir = Join-Path $HOME ".$AppName\bin"
$ConfigDir = Join-Path $HOME ".$AppName" # Standard directory for app data

Write-Host "--- Starting $AppName Uninstallation ---" -ForegroundColor Cyan

# 1. Kill Running Processes
# Prevent "Access Denied" errors by closing the app if it's currently running
$Process = Get-Process -Name $AppName -ErrorAction SilentlyContinue
if ($Process) {
    Write-Host "Closing running $AppName process..." -ForegroundColor Yellow
    Stop-Process -Name $AppName -Force
    Start-Sleep -Seconds 1 # Wait for file handles to release
}

# 2. Remove Binary and Bin Folder
if (Test-Path $InstallDir) {
    Write-Host "Removing binary folder: $InstallDir" -ForegroundColor Gray
    Remove-Item -Path $InstallDir -Recurse -Force -ErrorAction SilentlyContinue
}

# 3. Clean up Environment PATH
# This removes the entry we added during installation to keep the system clean
Write-Host "Cleaning up Environment PATH..." -ForegroundColor Gray
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -like "*$InstallDir*") {
    # Remove the path and any double semicolons resulting from the removal
    $NewPath = $UserPath -replace [regex]::Escape(";$InstallDir"), ""
    $NewPath = $NewPath -replace [regex]::Escape($InstallDir), ""
    [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
    $env:PATH = $NewPath # Update current session
    Write-Host "Removed $InstallDir from User PATH." -ForegroundColor Green
}

# 4. Optional: Data & Configuration Cleanup
$CleanupData = Read-Host "Delete all configuration files and vector databases? (y/n)"
if ($CleanupData -eq "y") {
    # 4a. Config File (.archon.yaml)
    $ConfigFile = Join-Path $HOME ".archon.yaml"
    if (Test-Path $ConfigFile) {
        Remove-Item -Path $ConfigFile -Force
        Write-Host "Deleted configuration file: $ConfigFile" -ForegroundColor Gray
    }

    # 4b. App Data Folder (chromem_db and others)
    if (Test-Path $ConfigDir) {
        Remove-Item -Path $ConfigDir -Recurse -Force
        Write-Host "Deleted app data directory: $ConfigDir" -ForegroundColor Gray
    }
}

Write-Host "`nSUCCESS: $AppName has been completely uninstalled." -ForegroundColor Green
Write-Host "Note: You may need to restart your terminal for PATH changes to fully sync." -ForegroundColor Yellow
