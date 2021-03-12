@echo off

set DIR=%~dp0%

rem Deleting packages
echo Deleting vagrant
echo
choco uninstall vagrant

echo Deleting virtualbox
echo
choco uninstall virtualbox

echo Deleting choco
rmdir C:\ProgramData\chocolatey