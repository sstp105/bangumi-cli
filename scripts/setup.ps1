# setup.ps1

$envVars = @{
    "LOCAL_SERVER_PORT"     = "8765"
    "BANGUMI_CLIENT_ID"     = "bangumi APP ID"
    "BANGUMI_CLIENT_SECRET" = "bangumi APP Secret"
    "QBITTORRENT_SERVER"    = "http://localhost:8767"
    "QBITTORRENT_USERNAME"  = "admin"
    "QBITTORRENT_PASSWORD"  = ""
    "MIKAN_IDENTITY_COOKIE" = ""
}

foreach ($key in $envVars.Keys) {
    Write-Host "Setting $key..."
    [Environment]::SetEnvironmentVariable($key, $envVars[$key], "Machine")
}

Write-Host "✅ Successfully added all envs."
