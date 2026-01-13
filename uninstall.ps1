<#
+--------------------------------------------------------------+
|  MODBUS MEMORY APPLIANCE (MMA) - WINDOWS UNINSTALLER          |
|  Service Removal via NSSM                                    |
|                                                              |
|  Gwapong Pilipino · Taas Noo · Kahit Kanino                   |
|                                                              |
|  Author : Rod Tamin                                          |
+--------------------------------------------------------------+
#
# Purpose:
#   Removes Modbus Memory Appliance (MMA) Windows service
#   installed via NSSM.
#
# Behavior:
#   - Stops service if running
#   - Removes NSSM service entry
#   - DOES NOT delete config.yaml by default
#   - DOES NOT touch data outside install directory
#
# Requirements:
#   - Run as Administrator
#>

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

Write-Host ""
Write-Host "+--------------------------------------------------------------+" -ForegroundColor Cyan
Write-Host "|  MODBUS MEMORY APPLIANCE (MMA) - WINDOWS UNINSTALLER          |" -ForegroundColor Cyan
Write-Host "|  Service Removal via NSSM                                    |" -ForegroundColor Cyan
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
$nssm        = Join-Path $PSScriptRoot "nssm.exe"

# ------------------------------------------------------------
# Admin check
# ------------------------------------------------------------
$identity  = [Security.Principal.WindowsIdentity]::GetCurrent()
$principal = New-Object Security.Principal.WindowsPrincipal($identity)

if (-not $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    throw "Uninstaller must be run as Administrator."
}

Write-Host "[OK] Administrator privileges confirmed."

# ------------------------------------------------------------
# Validate NSSM exists
# ------------------------------------------------------------
if (-not (Test-Path $nssm)) {
    throw "nssm.exe not found beside uninstall.ps1"
}

# ------------------------------------------------------------
# Stop service (if exists)
# ------------------------------------------------------------
Write-Host "[INFO] Stopping service (if running)..."

try {
    & $nssm stop $ServiceName | Out-Null
    Write-Host "[OK] Service stopped."
} catch {
    Write-Host "[WARN] Service not running or not present."
}

# ------------------------------------------------------------
# Remove service
# ------------------------------------------------------------
Write-Host "[INFO] Removing service registration..."

try {
    & $nssm remove $ServiceName confirm | Out-Null
    Write-Host "[OK] Service removed."
} catch {
    Write-Host "[WARN] Service not found in NSSM."
}

# ------------------------------------------------------------
# Install directory handling
# ------------------------------------------------------------
if (Test-Path $InstallDir) {
    Write-Host "[INFO] Install directory exists:"
    Write-Host "       $InstallDir"
    Write-Host "[INFO] Files were NOT deleted automatically."
    Write-Host "[INFO] You may remove this directory manually if desired."
} else {
    Write-Host "[OK] Install directory not present."
}

Write-Host ""
Write-Host "UNINSTALLATION COMPLETE." -ForegroundColor Green
Write-Host ""
