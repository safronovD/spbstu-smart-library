@echo off

set DIR=%~dp0%
mkdir DIR\logs
echo Starting application

cd DIR\..\
start vagrant up

set /p check="Print [exit] when you done: "
if "%check%" == "exit" (
    echo Cleaning data
    echo *************
    vagrant destroy
    if ERRORLEVEL 1 ( echo Cleaning error )
    exit 0
)

