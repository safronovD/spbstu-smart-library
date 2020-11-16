@echo off

SET DIR=%~dp0%

::setup elastic and kibana
%systemroot%\System32\WindowsPowerShell\v1.0\powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "& '%DIR%scripts\setup_docker.ps1' %*"