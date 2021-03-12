@echo off

set DIR=%~dp0%

echo Instalation started

rem chocolatey installtion
@"%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -NoProfile -InputFormat None -ExecutionPolicy Bypass -Command "[System.Net.ServicePointManager]::SecurityProtocol = 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))" && SET "PATH=%PATH%;%ALLUSERSPROFILE%\chocolatey\bin"

if ERRORLEVEL 1 ( 
    echo ***************************************************
    echo There is some troubles with chocolatey installation
)

rem packages installation
rem virtualbox
choco install virtualbox --version 6.1.18
if ERRORLEVEL 1 (
    echo ***************************************************
    echo There is some troubles with virtualbox installation
)

rem vagrant
choco install vagrant --version 2.2.14
if ERRORLEVEL 1 (
    echo ************************************************
    echo There is some troubles with vagrant installation
)

echo Some utilities needs reload. (Also you can do it manually later)
set /p check="Are you agree restart now [Yes]/[No]: "
set check=%check:~0,1%
if "%check%" == "Y" || "%check%" == "y" ( shutdown /r /t 0 )

exit /b 0