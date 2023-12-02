Write-Host "pre-applypatch hook called. Working dir is $(Get-Location)"
Write-Host "Environment:"
Get-ChildItem env: