<#
+--------------------------------------------------------------+
|  MODBUS MEMORY APPLIANCE (MMA) - WINDOWS INSTALLER            |
|  Service Setup via NSSM                                      |
|                                                              |
|  Gwapong Pilipino · Taas Noo · Kahit Kanino                   |
|                                                              |
|  Author : Rod Tamin                                          |
+--------------------------------------------------------------+
#
# Purpose:
#   Installs Modbus Memory Appliance (MMA) as a Windows service
#   using NSSM. This script performs NO configuration changes.
#
# Requirements:
#   - Run as Administrator
#   - mma.exe, nssm.exe, config.yaml in same directory
#
# Behavior:
#   - Fail fast
#   - Loud errors
#   - No silent fixes
#>

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

Write-Host ""
Write-Host "+--------------------------------------------------------------+" -ForegroundColor Cyan
Write-Host "|  MODBUS MEMORY APPLIANCE (MMA) - WINDOWS INSTALLER            |" -ForegroundColor Cyan
Write-Host "|  Service Setup via NSSM                                      |" -ForegroundColor Cyan
Write-Host "|                                                              |" -ForegroundColor Cyan
Write-Host "|  Gwapong Pilipino · Taas Noo · Kahit Kanino                   |" -ForegroundColor Cyan
Write-Host "|                                                              |" -ForegroundColor Cyan
Write-Host "|  Author : Rod Tamin                                          |" -ForegroundColor Cyan
Write-Host "+--------------------------------------------------------------+" -ForegroundColor Cyan
Write-Host ""

# ------------------------------------------------------------
# Constants
# ------------------------------------------------------------
$ServiceName = "ModbusMemoryAppliance"
$InstallDir  = "C:\Program Files\ModbusMemoryAppliance"
$HealthURL   = "http://localhost:8080/api/v1/health"

$RequiredFiles = @(
    "mma.exe",
    "nssm.exe",
    "config.yaml"
)

# ------------------------------------------------------------
# Admin check
# ------------------------------------------------------------
$identity  = [Security.Principal.WindowsIdentity]::GetCurrent()
$principal = New-Object Security.Principal.WindowsPrincipal($identity)

if (-not $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    throw "Installer must be run as Administrator."
}

Write-Host "[OK] Administrator privileges confirmed."

# ------------------------------------------------------------
# Validate required files
# ------------------------------------------------------------
foreach ($file in $RequiredFiles) {
    $path = Join-Path $PSScriptRoot $file
    if (-not (Test-Path $path)) {
        throw "Required file missing: $file (must be beside install.ps1)"
    }
}

Write-Host "[OK] Required files present."

# ------------------------------------------------------------
# Create install directory
# ------------------------------------------------------------
if (-not (Test-Path $InstallDir)) {
    New-Item -Path $InstallDir -ItemType Directory | Out-Null
    Write-Host "[OK] Created install directory: $InstallDir"
} else {
    Write-Host "[OK] Install directory already exists."
}

# ------------------------------------------------------------
# Copy binaries and config
# ------------------------------------------------------------
Copy-Item "$PSScriptRoot\mma.exe"     "$InstallDir\mma.exe"     -Force
Copy-Item "$PSScriptRoot\config.yaml" "$InstallDir\config.yaml" -Force

Write-Host "[OK] Files copied to $InstallDir"

# ------------------------------------------------------------
# NSSM service install
# ------------------------------------------------------------
$nssm = Join-Path $PSScriptRoot "nssm.exe"

# Clean up any existing service (idempotent)
& $nssm stop   $ServiceName 2>$null | Out-Null
& $nssm remove $ServiceName confirm 2>$null | Out-Null

# Install service
& $nssm install $ServiceName "$InstallDir\mma.exe" `
  "--config `"$InstallDir\config.yaml`""

& $nssm set $ServiceName Start SERVICE_AUTO_START

Write-Host "[OK] Service registered via NSSM."

# ------------------------------------------------------------
# Start service
# ------------------------------------------------------------
& $nssm start $ServiceName
Write-Host "[OK] Service started."

# ------------------------------------------------------------
# Health check (retry loop)
# ------------------------------------------------------------
Write-Host "[INFO] Waiting for service health..."

$healthy = $false
for ($i = 1; $i -le 10; $i++) {
    try {
        $resp = Invoke-RestMethod -Uri $HealthURL -TimeoutSec 2
        if ($resp.status -eq "ok") {
            $healthy = $true
            break
        }
    } catch {
        Start-Sleep -Seconds 1
    }
}

if (-not $healthy) {
    throw "Service failed health check. Verify config.yaml and logs."
}

Write-Host "[OK] Health check passed."
Write-Host ""
Write-Host "INSTALLATION SUCCESSFUL." -ForegroundColor Green
Write-Host ""
