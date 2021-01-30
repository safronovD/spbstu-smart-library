@echo off

SET DIR=%~dp0%

vagrant destroy

choco uninstall virtualbox vagrant
