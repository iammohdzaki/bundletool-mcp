$ErrorActionPreference = "Stop"
$InstallDir = Join-Path $HOME ".bundletool-mcp"
$BinDir = Join-Path $InstallDir "bin"
$ExeName = "bundletool-mcp-windows-amd64.exe"
$ExePath = Join-Path $BinDir "bundletool-mcp.exe"

Write-Host "Fetching latest release from GitHub..."
$Release = Invoke-RestMethod -Uri "https://api.github.com/repos/iammohdzaki/bundletool-mcp/releases/latest"
$Asset = $Release.assets | Where-Object { $_.name -eq $ExeName }

if (-not $Asset) {
    Write-Error "Could not find a Windows binary for the latest release."
    exit 1
}

if (-not (Test-Path $BinDir)) {
    New-Item -ItemType Directory -Force -Path $BinDir | Out-Null
}

Write-Host "Downloading $($Asset.name)..."
Invoke-WebRequest -Uri $Asset.browser_download_url -OutFile $ExePath

Write-Host "`nSuccessfully installed to $ExePath"

$EscapedPath = $ExePath -replace "\\", "\\"
Write-Host "`n======================================================="
Write-Host "Please add the following to your MCP configuration JSON:"
Write-Host "======================================================="
Write-Host "{
  `"mcpServers`": {
    `"bundletool`": {
      `"command`": `"$EscapedPath`",
      `"args`": []
    }
  }
}"
