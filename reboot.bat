
@echo off
color 1f
 taskkill /f /im MicEngine.exe
 start /d %cd% MicEngine.exe
 pause
