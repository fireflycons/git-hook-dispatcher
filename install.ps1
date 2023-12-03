$ErrorActionPreference = 'Stop'
$hookDir = ".githooks"
$hookDirAbs = Join-Path $env:USERPROFILE $hookDir

function Test-IsAdmin
{
    ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")
}

function Install-Release
{
    $r = Invoke-RestMethod -Uri https://api.github.com/repos/fireflycons/git-hook-dispatcher/releases/latest

    $zip = $r.assets |
        Where-Object { $_.content_type -eq "application/zip" -and $_.state -eq "uploaded" }

    if (-not $zip)
    {
        throw "Cannot find zip artifact in release"
    }

    try
    {
        Invoke-WebRequest -Uri $zip.browser_download_url -OutFile (Join-Path $env:TEMP git-hook-dispatcher.zip)

        if (-not (Test-Path -Path $hookDirAbs -PathType Container)) {
            New-Item -Path $env:USERPROFILE -ItemType Directory -Name $hookDir
        }

        $downloadedZip = Join-Path $env:TEMP git-hook-dispatcher.zip
        Expand-Archive -Path $downloadedZip -DestinationPath $hookDirAbs -Force
    }
    finally
    {
        if (Test-Path -Path $downloadedZip -PathType Leaf)
        {
            Remove-Item $downloadedZip
        }
    }
}

function Test-Symlink
{
    param
    (
        $Path
    )
    ((Get-Item $Path).Attributes.ToString() -match "ReparsePoint")
}

function Read-Html {

    param
    (
        [string]$htmlString
    )

    $bytes = [System.Text.Encoding]::Unicode.GetBytes($htmlString)
    $htmlFile = New-Object -Com 'HTMLFile'
    if ($htmlFile.PSObject.Methods.Name -Contains 'IHTMLDocument2_Write') {
        $htmlFile.IHTMLDocument2_Write($bytes)
    }
    else {
        $htmlFile.write($bytes)
    }
    $htmlFile.Close()
    $htmlFile
}

function Get-HookList
{
    @(
        'applypatch-msg'
        'pre-applypatch'
        'post-applypatch'
        'pre-commit'
        'pre-merge-commit'
        'prepare-commit-msg'
        'commit-msg'
        'post-commit'
        'pre-rebase'
        'post-checkout'
        'post-merge'
        'pre-push'
        'pre-receive'
        'update'
        'proc-receive'
        'post-receive'
        'post-update'
        'reference-transaction'
        'push-to-checkout'
        'pre-auto-gc'
        'post-rewrite'
        'sendemail-validate'
        'fsmonitor-watchman'
        'post-index-change'
    )
}

function Install-Symlinks
{
    Push-Location $hookDirAbs
    try
    {
        Get-HookList |
        ForEach-Object {
            $create = $false
            $link = "$_.exe"
            if (-not (Test-Path -Path $link -PathType Leaf))
            {
                $create = $true
            }
            elseif (-not (Test-Symlink -Path $link)) {
                # Remove any copied .exe
                Remove-Item $link
                $create = $true
            }

            if ($create)
            {
                New-Item -ItemType SymbolicLink -Name $link -Target hook.exe
            }
        }
    }
    finally
    {
        Pop-Location
    }
}

function Install-Copies
{

    Push-Location $hookDirAbs
    try
    {
        Get-HookList |
        ForEach-Object {
            Copy-Item -Path hook.exe -Destination "$_.exe"
        }
    }
    finally
    {
        Pop-Location
    }
}

Install-Release

if (Test-IsAdmin)
{
    Install-Symlinks $allHooks
}
else
{
    Install-Copies $allHooks
}

