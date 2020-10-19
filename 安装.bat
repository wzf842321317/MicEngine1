@echo off
cd /d %~dp0
call run-install.bat  MicEngine MicEngine.exe
call run-install.bat  MicDog micdog.exe
pause
