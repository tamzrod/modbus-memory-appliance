@echo off
echo +--------------------------------------------------------------+
echo |  MODBUS MEMORY APPLIANCE (MMA) - WINDOWS INSTALLER            |
echo |  Service Setup via NSSM                                      |
echo |                                                              |
echo |  Gwapong Pilipino · Taas Noo · Kahit Kanino                   |
echo |                                                              |
echo |  Author : Rod Tamin                                          |
echo +--------------------------------------------------------------+
echo.

powershell -NoProfile -ExecutionPolicy Bypass ^
  -File "%~dp0install.ps1"

IF ERRORLEVEL 1 (
  echo.
  echo INSTALLATION FAILED.
) ELSE (
  echo.
  echo INSTALLATION COMPLETED SUCCESSFULLY.
)

echo.
pause
