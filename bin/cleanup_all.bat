@echo off

SET DIR=%~dp0%

::cleanup
%systemroot%\System32\WindowsPowerShell\v1.0\powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "& '%DIR%scripts\cleanup.ps1' %*"