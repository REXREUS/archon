$AppName = "archon.exe"

Write-Host "Removing $AppName..." -ForegroundColor Cyan

$Commands = Get-Command $AppName -ErrorAction SilentlyContinue -All
if ($Commands) {
    foreach ($cmd in $Commands) {
        $path = $cmd.Source
        Write-Host "Removing $path"
        Remove-Item -Path $path -Force
    }
    Write-Host "$AppName has been removed." -ForegroundColor Green
} else {
    Write-Host "$AppName not found in PATH." -ForegroundColor Yellow
}

$DelConf = Read-Host "Delete configuration folder and vector database? (y/n)"
if ($DelConf -eq "y") {
    $ConfigFile = Join-Path $env:USERPROFILE ".archon.yaml"
    if (Test-Path $ConfigFile) {
        Remove-Item -Path $ConfigFile -Force
        Write-Host "Configuration file deleted."
    }

    if (Test-Path "chromem_db") {
        Remove-Item -Path "chromem_db" -Recurse -Force
        Write-Host "Vector database deleted."
    }
    
    Write-Host "Configuration and database have been deleted." -ForegroundColor Green
}

Write-Host "Uninstall complete." -ForegroundColor Cyan
