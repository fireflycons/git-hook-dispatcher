# reference-transaction hook needs to read from stdin
Write-Host "refrence-transaction started"

$inputStream = [System.Console]::In

while ($line = $inputStream.ReadLine())
{
    Write-Host "reference-transaction received stdin:" $line
}

