@echo off
echo ===============================================
echo Modbus Memory Appliance - Windows Uninstaller
echo ===============================================
echo.

powershell -NoProfile -ExecutionPolicy Bypass ^
  -File "%~dp0uninstall.ps1"

IF ERRORLEVEL 1 (
  echo.
  echo UNINSTALL FAILED.
) ELSE (
  echo.
  echo UNINSTALL COMPLETED.
)

echo.
pause
